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

// This file defines file manipulation operations.

package buflsp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/bufbuild/buf/private/buf/bufworkspace"
	"github.com/bufbuild/buf/private/bufpkg/bufanalysis"
	"github.com/bufbuild/buf/private/bufpkg/bufcheck"
	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/pkg/slogext"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/protocompile/ast"
	"go.lsp.dev/protocol"
)

// fileold is a fileold that has been opened by the client.
//
// Mutating a fileold is thread-safe.
type fileold struct {
	lsp *lsp
	uri protocol.URI

	text string
	// Version is an opaque version identifier given to us by the LSP client. This
	// is used in the protocol to disambiguate which version of a file e.g. publishing
	// diagnostics or symbols an operating refers to.
	version int32
	hasText bool // Whether this file has ever had text read into it.
	// Always set false->true. Once true, never becomes false again.

	workspace bufworkspace.Workspace
	module    bufmodule.Module

	objectInfo storage.ObjectInfo
	fileInfo   bufmodule.FileInfo

	fileNode    *ast.FileNode
	diagnostics []protocol.Diagnostic
	symbols     []*symbol
	image       bufimage.Image
}

// IsWKT returns whether this file corresponds to a well-known type.
func (f *fileold) IsWKT() bool {
	_, ok := f.objectInfo.(wktObjectInfo)
	return ok
}

// IsLocal returns whether this is a local file, i.e. a file that the editor
// is editing and not something from e.g. the BSR.
func (f *fileold) IsLocal() bool {
	if f.objectInfo == nil {
		return false
	}

	return f.objectInfo.LocalPath() == f.objectInfo.ExternalPath()
}

// Package returns the package of this file, if known.
func (f *fileold) Package() string {
	var pkg string
	if f.fileNode != nil {
		ast.Walk(f.fileNode, &ast.SimpleVisitor{
			DoVisitPackageNode: func(pn *ast.PackageNode) error {
				pkg = string(pn.Name.AsIdentifier())
				return nil
			},
		})
	}

	return pkg
}

// Refresh rebuilds all of a file's internal book-keeping.
//
// If deep is set, this will also load imports and refresh those, too.
func (f *fileold) Refresh(ctx context.Context) {
	// f.RunLints(ctx)
}

// RunLints runs linting on this file. Returns whether any lints failed.
//
// This operation requires BuildImage().
func (f *fileold) RunLints(ctx context.Context) bool {
	if f.IsWKT() {
		// Well-known types are not linted.
		return false
	}

	workspace := f.workspace
	module := f.module
	image := f.image

	if module == nil || image == nil {
		f.lsp.logger.Warn(fmt.Sprintf("could not find image for %q", f.uri))
		return false
	}

	f.lsp.logger.Debug(fmt.Sprintf("running lint for %q in %v", f.uri, module.ModuleFullName()))

	lintConfig := workspace.GetLintConfigForOpaqueID(module.OpaqueID())
	err := f.lsp.checkClient.Lint(
		ctx,
		lintConfig,
		image,
		bufcheck.WithPluginConfigs(workspace.PluginConfigs()...),
	)

	if err == nil {
		f.lsp.logger.Warn(fmt.Sprintf("lint generated no errors for %s", f.uri))
		return false
	}

	var annotations bufanalysis.FileAnnotationSet
	if !errors.As(err, &annotations) {
		f.lsp.logger.Warn("error while linting", slog.String("uri", string(f.uri)), slogext.ErrorAttr(err))
		return false
	}

	f.lsp.logger.Warn(fmt.Sprintf("lint generated %d error(s) for %s", len(annotations.FileAnnotations()), f.uri))

	for _, annotation := range annotations.FileAnnotations() {
		f.lsp.logger.Info(annotation.FileInfo().Path(), " ", annotation.FileInfo().ExternalPath())

		f.diagnostics = append(f.diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(annotation.StartLine()) - 1,
					Character: uint32(annotation.StartColumn()) - 1,
				},
				End: protocol.Position{
					Line:      uint32(annotation.EndLine()) - 1,
					Character: uint32(annotation.EndColumn()) - 1,
				},
			},
			Code:     annotation.Type(),
			Severity: protocol.DiagnosticSeverityError,
			Source:   LSPServerName,
			Message:  annotation.Message(),
		})
	}
	return true
}

// SymbolAt finds a symbol in this file at the given cursor position, if one exists.
//
// Returns nil if no symbol is found.
func (f *fileold) SymbolAt(ctx context.Context, cursor protocol.Position) *symbol {
	// Binary search for the symbol whose start is before or equal to cursor.
	idx, found := slices.BinarySearchFunc(f.symbols, cursor, func(sym *symbol, cursor protocol.Position) int {
		return comparePositions(sym.Range().Start, cursor)
	})
	if !found {
		if idx == 0 {
			return nil
		}
		idx--
	}

	symbol := f.symbols[idx]

	// Check that cursor is before the end of the symbol.
	if comparePositions(symbol.Range().End, cursor) <= 0 {
		return nil
	}

	return symbol
}

// wktObjectInfo is a concrete type to help us identify WKTs among the
// importable files.
type wktObjectInfo struct {
	storage.ObjectInfo
}
