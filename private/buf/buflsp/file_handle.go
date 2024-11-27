package buflsp

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bufbuild/buf/private/buf/buflsp/file"
	"github.com/bufbuild/buf/private/buf/buflsp/parser"
	"github.com/bufbuild/protocompile/ast"
	"go.lsp.dev/protocol"
)

type fileHandle struct {
	handle     file.Handle
	parsedFile *parser.File
}

func (fh *fileHandle) Format() (string, protocol.Range, error) {
	var out strings.Builder
	if err := fh.parsedFile.Format(&out); err != nil {
		return "", protocol.Range{}, fmt.Errorf("parser file format: %w", err)
	}

	return out.String(), fh.parsedFile.Range(), nil
}

func (fh *fileHandle) NodeAtPosition(pos protocol.Position) ast.Node {
	return fh.parsedFile.Nearest(pos)
}

func (fh *fileHandle) node2tok(n ast.Node, typ SemanticTokenType, ms ...SemanticTokenModifier) SemanticToken {
	ni := fh.parsedFile.NodeInfo(n)
	start := parser.NewZeroBaseSourcePos(ni.Start()).ToPosition()
	raw := ni.RawText()

	return SemanticToken{
		Line:      start.Line,
		Start:     start.Character,
		Len:       uint32(len(raw)),
		Raw:       raw,
		Type:      typ,
		Modifiers: ms,
	}
}

func (fh *fileHandle) SemanticTokens() []SemanticToken {
	toks := []SemanticToken{}

	_ = fh.parsedFile.Walk(&ast.SimpleVisitor{
		DoVisitKeywordNode: func(n *ast.KeywordNode) error {
			toks = append(toks, fh.node2tok(n, SemanticTokenTypeKeyWord))
			return nil
		},
		DoVisitPackageNode: func(n *ast.PackageNode) error {
			toks = append(toks, fh.node2tok(n.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(n.Name, SemanticTokenTypeNamespace, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitServiceNode: func(n *ast.ServiceNode) error {
			toks = append(toks, fh.node2tok(n.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(n.Name, SemanticTokenTypeInterface, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitRPCNode: func(val *ast.RPCNode) error {
			toks = append(toks, fh.node2tok(val.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeMethod, SemanticTokenModifierDeclaration))

			if val.Input.Stream != nil {
				toks = append(toks, fh.node2tok(val.Input.Stream, SemanticTokenTypeKeyWord))
			}
			if val.Input.MessageType != nil {
				toks = append(toks, fh.node2tok(val.Input.MessageType, SemanticTokenTypeType))
			}

			toks = append(toks, fh.node2tok(val.Returns, SemanticTokenTypeKeyWord))

			if val.Output.Stream != nil {
				toks = append(toks, fh.node2tok(val.Output.Stream, SemanticTokenTypeKeyWord))
			}
			if val.Output.MessageType != nil {
				toks = append(toks, fh.node2tok(val.Output.MessageType, SemanticTokenTypeType))
			}
			return nil
		},
		DoVisitMessageNode: func(val *ast.MessageNode) error {
			toks = append(toks, fh.node2tok(val.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeStruct, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitEnumNode: func(val *ast.EnumNode) error {
			toks = append(toks, fh.node2tok(val.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeEnum, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitEnumValueNode: func(val *ast.EnumValueNode) error {
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeEnumMember, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitFieldNode: func(val *ast.FieldNode) error {
			if val.Label.KeywordNode != nil {
				toks = append(toks, fh.node2tok(val.Label.KeywordNode, SemanticTokenTypeKeyWord))
			}
			toks = append(toks, fh.node2tok(val.FldType, SemanticTokenTypeType, SemanticTokenModifierDeclaration))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeProperty, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitMapFieldNode: func(val *ast.MapFieldNode) error {
			toks = append(toks, fh.node2tok(val.MapType.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(val.MapType.KeyType, SemanticTokenTypeType, SemanticTokenModifierDeclaration))
			toks = append(toks, fh.node2tok(val.MapType.ValueType, SemanticTokenTypeType, SemanticTokenModifierDeclaration))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeProperty, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitOneofNode: func(val *ast.OneofNode) error {
			toks = append(toks, fh.node2tok(val.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeStruct, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitGroupNode: func(val *ast.GroupNode) error {
			if val.Label.KeywordNode != nil {
				toks = append(toks, fh.node2tok(val.Label.KeywordNode, SemanticTokenTypeKeyWord))
			}
			toks = append(toks, fh.node2tok(val.Keyword, SemanticTokenTypeKeyWord))
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeStruct, SemanticTokenModifierDeclaration))
			return nil
		},
		DoVisitOptionNode: func(val *ast.OptionNode) error {
			if val.Keyword != nil {
				toks = append(toks, fh.node2tok(val.Keyword, SemanticTokenTypeKeyWord))
			}
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeDecorator, SemanticTokenModifierDefinition))
			return nil
		},
		DoVisitFieldReferenceNode: func(val *ast.FieldReferenceNode) error {
			toks = append(toks, fh.node2tok(val.Name, SemanticTokenTypeDecorator, SemanticTokenModifierDefinition))
			return nil
		},
		DoVisitStringValueNode: func(val ast.StringValueNode) error {
			toks = append(toks, fh.node2tok(val, SemanticTokenTypeString))
			return nil
		},
		DoVisitIntValueNode: func(val ast.IntValueNode) error {
			toks = append(toks, fh.node2tok(val, SemanticTokenTypeNumber))
			return nil
		},
		DoVisitFloatValueNode: func(val ast.FloatValueNode) error {
			toks = append(toks, fh.node2tok(val, SemanticTokenTypeNumber))
			return nil
		},
		DoVisitSignedFloatLiteralNode: func(val *ast.SignedFloatLiteralNode) error {
			toks = append(toks, fh.node2tok(val, SemanticTokenTypeNumber))
			return nil
		},
	})

	slices.SortFunc(toks, func(left, right SemanticToken) int {
		if left.Line != right.Line {
			return int(left.Line - right.Line)
		} else {
			return int(left.Start - right.Start)
		}
	})

	return toks
}
