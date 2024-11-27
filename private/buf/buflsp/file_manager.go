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

// This file defines a manager for tracking individual files.

package buflsp

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/buf/private/buf/buflsp/parser"
	"github.com/bufbuild/buf/private/pkg/slogext"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/protocompile/ast"
	"go.lsp.dev/protocol"
)

// fileManager tracks all files the LSP is currently handling, whether read from disk or opened
// by the editor.
type fileManager struct {
	lsp    *lsp
	logger *slog.Logger

	fs     *filesystem
	parser *parser.Parser

	diagnosticClient *DiagnosticClient
}

// newFiles creates a new file manager.
func newFileManager(lsp *lsp) *fileManager {

	dc := NewDiagnosticClient(lsp.logger, lsp.client)
	fs := NewFileSystem(lsp.logger, lsp.controller)
	psr := parser.NewParser(lsp.logger, dc, fs.Resolver())

	return &fileManager{
		lsp:              lsp,
		logger:           lsp.logger,
		fs:               fs,
		parser:           psr,
		diagnosticClient: dc,
	}
}

func (fm *fileManager) Init(ctx context.Context, folders protocol.WorkspaceFolders) error {
	return fm.fs.init(ctx, folders)
}

func (fm *fileManager) Open(ctx context.Context, overlay *file.Overlay) (*fileHandle, error) {
	fm.logger.DebugContext(
		ctx, "open overlay",
		"Overlay", overlay,
	)

	fh, err := fm.fs.Open(overlay)
	if err != nil {
		return nil, fmt.Errorf("fs open: %w", err)
	}

	return fm.ReadFile(ctx, fh.URI())
}

func (fm *fileManager) Change(ctx context.Context, overlay *file.Overlay) (*fileHandle, error) {
	fh, err := fm.fs.Open(overlay)
	if err != nil {
		return nil, fmt.Errorf("fs open: %w", err)
	}

	stat, err := fm.fs.Source().Stat(ctx, fh.URI())
	if err != nil {
		return nil, fmt.Errorf("fs stat file: %w", err)
	}

	return fm.internalparse(ctx, stat, fh)
}

func (fm *fileManager) Get(ctx context.Context, uri protocol.DocumentURI) (*fileHandle, error) {
	stat, err := fm.fs.Source().Stat(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("fs stat file: %w", err)
	}

	fh, err := fm.fs.Source().ReadFile(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("fs read file: %w", err)
	}

	pf, err := fm.parser.Parse(file.NormalURIStr(stat.Path()), fh)
	if err != nil {
		return nil, fmt.Errorf("parser parse file: %w", err)
	}

	return &fileHandle{
		handle:     fh,
		parsedFile: pf,
	}, nil
}

// ReadFile will recursively read from fs and parse all imports
func (fm *fileManager) ReadFile(ctx context.Context, uri protocol.DocumentURI) (*fileHandle, error) {
	fm.logger.DebugContext(ctx, "read file", "URI", uri)

	stat, err := fm.fs.Source().Stat(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("fs stat file: %w", err)
	}

	fh, err := fm.fs.Source().ReadFile(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("fs read file: %w", err)
	}

	return fm.internalparse(ctx, stat, fh)
}

func (fm *fileManager) internalparse(ctx context.Context, stat storage.ObjectInfo, handle file.Handle) (*fileHandle, error) {
	defer slogext.Profile(fm.logger, "URI", handle.URI())()

	fm.diagnosticClient.Reset(handle)

	pf, err := fm.parser.Parse(file.NormalURIStr(stat.Path()), handle)
	if err != nil {
		return nil, fmt.Errorf("parser parse file: %w", err)
	}

	if err := fm.diagnosticClient.Notify(handle); err != nil {
		fm.logger.WarnContext(
			ctx, "notify diagnostic fail",
			"URI", handle.URI(),
			"Error", err,
		)
	}

	for _, imp := range pf.Imports() {
		if _, err := fm.ReadFile(ctx, imp); err != nil {
			return nil, err
		}
	}

	return &fileHandle{
		handle:     handle,
		parsedFile: pf,
	}, nil
}

// Close marks a file as closed.
//
// This will not necessarily evict the file, since there may be more than one user
// for this file.
func (fm *fileManager) Close(ctx context.Context, uri protocol.URI) {
	fm.fs.Source().Close(uri)
}

func (fm *fileManager) Definition(ctx context.Context, fh *fileHandle, pos protocol.Position) ([]protocol.Location, error) {
	defer slogext.Profile(fm.logger, "File", fh.handle.URI(), "Position", pos)()

	if fh == nil {
		return nil, os.ErrNotExist
	}

	ns := fh.parsedFile.NodesAt(pos)
	if len(ns) == 0 {
		return nil, nil // no nodes, no definitions
	}

	// two cases for now, import or ref
	imports := []*ast.ImportNode{}

	for _, n := range ns {
		switch val := n.(type) {
		case *ast.ImportNode:
			imports = append(imports, val)
		}
	}

	locs := []protocol.Location{}

	// only one
	if len(imports) > 0 {
		impname := imports[0].Name.AsString()
		locs = append(locs, fm.fs.Location(ctx, file.NormalURIStr(impname))...)
		fm.logger.DebugContext(ctx, "import location", "Import", impname, "Locations", locs)
	}

	return locs, nil
}
