package parser

import (
	"log/slog"

	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/protocompile/reporter"
	"go.lsp.dev/protocol"
)

var source = "buf-lsp"

var noopReporter = reporter.NewReporter(
	func(err reporter.ErrorWithPos) error {
		return nil
	},
	func(ewp reporter.ErrorWithPos) {},
)

type fileReporter struct {
	logger              *slog.Logger
	diagnosticCollector DiagnosticCollector
	handle              file.Handle
}

func (p *fileReporter) Error(err reporter.ErrorWithPos) error {
	p.logger.Debug(
		"report error",
		"Start", err.Start(),
		"End", err.End(),
		"Position", err.GetPosition(),
		"ErrMsg", err.Error(),
	)

	start, end := NewZeroBaseSourcePos(err.Start()), NewZeroBaseSourcePos(err.End())

	_ = p.diagnosticCollector.AddDiagnostics(p.handle, []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: start.ToPosition(),
				End:   end.ToPosition(),
			},
			Severity: protocol.DiagnosticSeverityError,
			Message:  err.Unwrap().Error(),
			Source:   source,
		},
	})
	return nil
}

func (p *fileReporter) Warning(warn reporter.ErrorWithPos) {
	p.logger.Debug(
		"report warning",
		"Start", warn.Start(),
		"End", warn.End(),
		"Position", warn.GetPosition(),
		"WarnMsg", warn.Error(),
	)

	start, end := NewZeroBaseSourcePos(warn.Start()), NewZeroBaseSourcePos(warn.End())

	_ = p.diagnosticCollector.AddDiagnostics(p.handle, []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: start.ToPosition(),
				End:   end.ToPosition(),
			},
			Severity: protocol.DiagnosticSeverityWarning,
			Message:  warn.Unwrap().Error(),
			Source:   source,
		},
	})
}
