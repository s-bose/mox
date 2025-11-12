package ast

import (
	"testing"

	"github.com/s-bose/mox/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "foo"}, Value: "foo"},
				Type:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "int"}, Value: "int"},
				Value: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "bar"}, Value: "bar"},
			},
		},
	}

	if program.String() != "let foo: int = bar;" {
		t.Errorf("program.String() invalid, got=%q", program.String())
	}
}
