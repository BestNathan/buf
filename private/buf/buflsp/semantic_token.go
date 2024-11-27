package buflsp

import (
	"go.lsp.dev/protocol"
)

const SemanticTokenDecorator protocol.SemanticTokenTypes = "decorator"

const (
	SemanticTokenTypeNamespace SemanticTokenType = iota
	SemanticTokenTypeType
	SemanticTokenTypeEnum
	SemanticTokenTypeStruct
	SemanticTokenTypeInterface
	SemanticTokenTypeTypeParameter
	SemanticTokenTypeParameter
	SemanticTokenTypeProperty
	SemanticTokenTypeEnumMember
	SemanticTokenTypeMethod
	SemanticTokenTypeKeyWord
	SemanticTokenTypeModifier
	SemanticTokenTypeString
	SemanticTokenTypeNumber
	SemanticTokenTypeDecorator
	SemanticTokenTypeLastButNoMeaning // just for range
)

const (
	SemanticTokenModifierDeclaration SemanticTokenModifier = 1 << iota
	SemanticTokenModifierDefinition
	SemanticTokenModifierLastButNoMeaning
)

func SematicTokenTypes() []protocol.SemanticTokenTypes {
	typs := []protocol.SemanticTokenTypes{}
	for idx := range SemanticTokenTypeLastButNoMeaning {
		if typ, ok := SemanticTokenTypeMapping[idx]; ok {
			typs = append(typs, typ)
		}
	}
	return typs
}

func SematicTokenModifiers() []protocol.SemanticTokenModifiers {
	mdfs := []protocol.SemanticTokenModifiers{}
	for idx := range SemanticTokenModifierLastButNoMeaning {
		if mdf, ok := SemanticTokenModifierMapping[idx]; ok {
			mdfs = append(mdfs, mdf)
		}
	}
	return mdfs
}

var (
	SemanticTokenTypeMapping = map[SemanticTokenType]protocol.SemanticTokenTypes{
		SemanticTokenTypeNamespace:     protocol.SemanticTokenNamespace,
		SemanticTokenTypeType:          protocol.SemanticTokenType,
		SemanticTokenTypeEnum:          protocol.SemanticTokenEnum,
		SemanticTokenTypeStruct:        protocol.SemanticTokenStruct,
		SemanticTokenTypeInterface:     protocol.SemanticTokenInterface,
		SemanticTokenTypeTypeParameter: protocol.SemanticTokenTypeParameter,
		SemanticTokenTypeParameter:     protocol.SemanticTokenParameter,
		SemanticTokenTypeProperty:      protocol.SemanticTokenProperty,
		SemanticTokenTypeEnumMember:    protocol.SemanticTokenEnumMember,
		SemanticTokenTypeMethod:        protocol.SemanticTokenMethod,
		SemanticTokenTypeKeyWord:       protocol.SemanticTokenKeyword,
		SemanticTokenTypeModifier:      protocol.SemanticTokenModifier,
		SemanticTokenTypeString:        protocol.SemanticTokenString,
		SemanticTokenTypeNumber:        protocol.SemanticTokenNumber,
		SemanticTokenTypeDecorator:     SemanticTokenDecorator,
	}
	SemanticTokenModifierMapping = map[SemanticTokenModifier]protocol.SemanticTokenModifiers{
		SemanticTokenModifierDeclaration: protocol.SemanticTokenModifierDeclaration,
		SemanticTokenModifierDefinition:  protocol.SemanticTokenModifierDefinition,
	}
	SemanticTokensLegend = protocol.SemanticTokensLegend{
		TokenTypes:     SematicTokenTypes(),
		TokenModifiers: SematicTokenModifiers(),
	}
)

type SemanticTokenType uint32

type SemanticTokenModifier uint32
type SemanticTokenModifiers []SemanticTokenModifier

func (ms SemanticTokenModifiers) Value() uint32 {
	if len(ms) == 0 {
		return 0
	}

	var val uint32

	for _, v := range ms {
		val |= uint32(v)
	}

	return val
}

type SemanticToken struct {
	Line, Start uint32
	Len         uint32
	Raw         string
	Type        SemanticTokenType
	Modifiers   SemanticTokenModifiers
}

func EncodeSemanticTokens(toks []SemanticToken) []uint32 {
	// each semantic token needs five values
	// (see Integer Encoding for Tokens in the LSP spec)
	x := make([]uint32, 5*len(toks))
	var j int
	var last SemanticToken
	for i := 0; i < len(toks); i++ {
		item := toks[i]

		if j == 0 {
			x[0] = toks[0].Line
		} else {
			x[j] = item.Line - last.Line
		}
		x[j+1] = item.Start
		if j > 0 && x[j] == 0 {
			x[j+1] = item.Start - last.Start
		}
		x[j+2] = item.Len
		x[j+3] = uint32(item.Type)
		x[j+4] = item.Modifiers.Value()
		j += 5
		last = item
	}
	return x[:j]
}
