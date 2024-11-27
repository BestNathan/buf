package buflsp

import (
	"fmt"
	"testing"
)

func TestSemanticTokenModifer(t *testing.T) {
	ms := SemanticTokenModifiers{SemanticTokenModifierDeclaration, SemanticTokenModifierDefinition}
	t.Log(ms.Value())
	t.Log(fmt.Sprintf("0x%b", ms.Value()))
}
