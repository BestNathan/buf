package buflsp

import "go.lsp.dev/protocol"

type SemanticTokensOptions struct {
	protocol.WorkDoneProgressOptions

	Legend protocol.SemanticTokensLegend `json:"legend"`
	Full   bool                          `json:"full"`
}
