package parser

import (
	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"go.lsp.dev/protocol"
)

type DiagnosticCollector interface {
	Reset(file.Handle)
	AddDiagnostics(file.Handle, []protocol.Diagnostic) error
}

type noopDiagnosticCollector struct{}

func (n *noopDiagnosticCollector) Reset(_ file.Handle) {
}

func (n *noopDiagnosticCollector) AddDiagnostics(_ file.Handle, _ []protocol.Diagnostic) error {
	return nil
}
