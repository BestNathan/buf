package parser

import (
	"github.com/bufbuild/protocompile/ast"
	"go.lsp.dev/protocol"
)

type PositionType int

const (
	PositionTypeImport = PositionType(iota)
	PositionTypeDecleration
	PositionTypeReference
)

type OneBasePosition protocol.Position

func NewOneBasePosition(pos protocol.Position) OneBasePosition {
	pos.Line++
	pos.Character++
	return OneBasePosition(pos)
}

type ZeroBaseSourcePos ast.SourcePos

func NewZeroBaseSourcePos(pos ast.SourcePos) ZeroBaseSourcePos {
	if pos.Col > 0 && pos.Line > 0 {
		pos.Line--
		pos.Col--
	}
	return ZeroBaseSourcePos(pos)
}

func (pos ZeroBaseSourcePos) ToPosition() protocol.Position {
	return protocol.Position{
		Line:      uint32(pos.Line),
		Character: uint32(pos.Col),
	}
}

func PositionWithinNode(ni ast.SourceSpan, zpos protocol.Position) bool {
	pos := NewOneBasePosition(zpos)

	start, end := ni.Start(), ni.End()

	// one line node, need check col
	if start.Line == end.Line && start.Line == int(pos.Line) {
		return start.Col <= int(pos.Character) && end.Col >= int(pos.Character)
	} else {
		return start.Line <= int(pos.Line) && end.Line >= int(pos.Line)
	}
}

func PositionNodeLen(ni ast.SourceSpan, zpos protocol.Position) (linelen int, collen int) {
	linelen, collen = -1, -1
	if !PositionWithinNode(ni, zpos) {
		return
	}

	pos := NewOneBasePosition(zpos)

	start, end := ni.Start(), ni.End()

	if start.Line == end.Line {
		return 0, (int(pos.Character) - start.Col) + (end.Col - int(pos.Character))
	} else {
		return (int(pos.Line) - start.Line) + (end.Line - int(pos.Line)), -1
	}
}
