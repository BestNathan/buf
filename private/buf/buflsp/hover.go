package buflsp

import "go.lsp.dev/protocol"

type HoverValuer interface {
	HoverValue() protocol.Hover
}

