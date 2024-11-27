// Copyright 2020-2024 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package buflsp implements a language server for Protobuf.
//
// The main entry-point of this package is the Serve() function, which creates a new LSP server.
package buflsp

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"

	internalfile "github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/protocompile/ast"
	"go.lsp.dev/protocol"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	LSPServerName = "buf-lsp"
)

// server is an implementation of protocol.Server.
//
// This is a separate type from buflsp.lsp so that the dozens of handler methods for this
// type are kept separate from the rest of the logic.
//
// See https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification.
type server struct {
	// This automatically implements all of protocol.Server for us. By default,
	// every method returns an error.
	nyi

	// We embed the LSP pointer as well, since it only has private members.
	*lsp
}

// newServer creates a protocol.Server implementation out of an lsp.
func newServer(lsp *lsp) protocol.Server {
	return &server{lsp: lsp}
}

// Methods for server are grouped according to the groups in the LSP protocol specification.

// -- Lifecycle Methods

// Initialize is the first message the LSP receives from the client. This is where all
// initialization of the server wrt to the project is is invoked on must occur.
func (s *server) Initialize(
	ctx context.Context,
	params *protocol.InitializeParams,
) (*protocol.InitializeResult, error) {
	if err := s.init(ctx, params); err != nil {
		return nil, err
	}

	if err := s.fileManager.Init(ctx, params.WorkspaceFolders); err != nil {
		return nil, fmt.Errorf("file manager init: %w", err)
	}

	info := &protocol.ServerInfo{Name: LSPServerName}
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		info.Version = buildInfo.Main.Version
	}

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			// These are all the things we advertise to the client we can do.
			// For now, incomplete features are explicitly disabled here as TODOs.
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				// Request that whole files be sent to us. Protobuf IDL files don't
				// usually get especially huge, so this simplifies our logic without
				// necessarily making the LSP slow.
				Change: protocol.TextDocumentSyncKindFull,
			},
			DefinitionProvider: &protocol.DefinitionOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{WorkDoneProgress: true},
			},
			DocumentFormattingProvider: true,
			HoverProvider:              true,
			SemanticTokensProvider: &SemanticTokensOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{WorkDoneProgress: true},
				Legend:                  SemanticTokensLegend,
				Full:                    true,
			},
		},
		ServerInfo: info,
	}, nil
}

// Initialized is sent by the client after it receives the Initialize response and has
// initialized itself. This is only a notification.
func (s *server) Initialized(
	ctx context.Context,
	params *protocol.InitializedParams,
) error {
	return nil
}

func (s *server) SetTrace(
	ctx context.Context,
	params *protocol.SetTraceParams,
) error {
	s.lsp.traceValue.Store(&params.Value)
	return nil
}

// Shutdown is sent by the client when it wants the server to shut down and exit.
// The client will wait until Shutdown returns, and then call Exit.
func (s *server) Shutdown(ctx context.Context) error {
	return nil
}

// Exit is a notification that the client has seen shutdown complete, and that the
// server should now exit.
func (s *server) Exit(ctx context.Context) error {
	// TODO: return an error if Shutdown() has not been called yet.

	// Close the connection. This will let the server shut down gracefully once this
	// notification is replied to.
	return s.lsp.conn.Close()
}

// -- File synchronization methods.

// DidOpen is called whenever the client opens a document. This is our signal to parse
// the file.
func (s *server) DidOpen(
	ctx context.Context,
	params *protocol.DidOpenTextDocumentParams,
) error {
	ol := internalfile.NewOverlay(
		params.TextDocument.URI,
		[]byte(params.TextDocument.Text),
		params.TextDocument.Version,
	)

	if _, err := s.fileManager.Open(ctx, ol); err != nil {
		return fmt.Errorf("file manager open: %w", err)
	}

	return nil
}

// DidChange is called whenever the client opens a document. This is our signal to parse
// the file.
func (s *server) DidChange(
	ctx context.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	if len(params.ContentChanges) == 0 {
		return nil
	}

	ol := internalfile.NewOverlay(
		params.TextDocument.URI,
		[]byte(params.ContentChanges[0].Text),
		params.TextDocument.Version,
	)

	_, err := s.fileManager.Change(ctx, ol)
	if err != nil {
		return fmt.Errorf("file manager change: %w", err)
	}

	return nil
}

// Formatting is called whenever the user explicitly requests formatting.
func (s *server) Formatting(
	ctx context.Context,
	params *protocol.DocumentFormattingParams,
) ([]protocol.TextEdit, error) {
	fh, err := s.fileManager.Get(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("file manager get: %w", err)
	}

	newtext, r, err := fh.Format()
	if err != nil {
		return nil, fmt.Errorf("file handle format: %w", err)
	}

	return []protocol.TextEdit{
		{NewText: newtext, Range: r},
	}, nil
}

// DidOpen is called whenever the client opens a document. This is our signal to parse
// the file.
func (s *server) DidClose(
	ctx context.Context,
	params *protocol.DidCloseTextDocumentParams,
) error {
	s.fileManager.Close(ctx, params.TextDocument.URI)
	return nil
}

// -- Language functionality methods.

// Hover is the entry point for hover inlays.
func (s *server) Hover(
	ctx context.Context,
	params *protocol.HoverParams,
) (*protocol.Hover, error) {
	fh, err := s.fileManager.Get(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("file manager get: %w", err)
	}

	node := fh.NodeAtPosition(params.Position)
	var nodespan ast.SourceSpan
	var fullname string
	if inode, ok := node.(*ast.IdentNode); ok {
		fullname = fmt.Sprintf("%s.%s", fh.parsedFile.PackageName(), inode.AsIdentifier())
		nodespan = fh.parsedFile.LookUpSymbol(protoreflect.FullName(fullname))
	}

	s.logger.Debug(
		"node at position",
		"Node", node,
		"NodeType", reflect.TypeOf(node),
		"Position", params.Position,
		"FullName", fullname,
		"Span", fmt.Sprintf("%#v", nodespan),
	)

	nodes := fh.parsedFile.NodesAt(params.Position)
	s.logger.Debug("nodes at position", "NodesType", nodes.TypeNames(), "Position", params.Position)

	return nil, nil
}

// Definition is the entry point for go-to-definition.
func (s *server) Definition(
	ctx context.Context,
	params *protocol.DefinitionParams,
) ([]protocol.Location, error) {
	fh, err := s.fileManager.Get(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("file manager get: %w", err)
	}

	progress := newProgressFromClient(s.lsp, &params.WorkDoneProgressParams)
	progress.Begin(ctx, "Searching")
	defer progress.Done(ctx)

	locs, err := s.fileManager.Definition(ctx, fh, params.Position)
	if err != nil {
		s.logger.Debug(
			"test fm definition fail",
			"Error", err,
		)
	}

	if len(locs) > 0 {
		return locs, nil
	}

	// symbol := file.SymbolAt(ctx, params.Position)
	// if symbol == nil {
	// 	return nil, nil
	// }

	// if _, ok := symbol.kind.(*import_); ok {
	// 	// This is an import, we just want to jump to the file.
	// 	// return []protocol.Location{{URI: imp.file.uri}}, nil
	// 	return nil, nil
	// }

	// def, _ := symbol.Definition(ctx)
	// if def != nil {
	// 	return []protocol.Location{{
	// 		URI:   def.file.uri,
	// 		Range: def.Range(),
	// 	}}, nil
	// }

	return nil, nil
}

// SemanticTokensFull is called to render semantic token information on the client.
func (s *server) SemanticTokensFull(
	ctx context.Context,
	params *protocol.SemanticTokensParams,
) (*protocol.SemanticTokens, error) {
	fh, err := s.fileManager.Get(ctx, params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("file manager get: %w", err)
	}

	toks := fh.SemanticTokens()

	return &protocol.SemanticTokens{
		Data: EncodeSemanticTokens(toks),
	}, nil
}
