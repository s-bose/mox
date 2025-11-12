package ast

import (
	"testing"

	"github.com/s-bose/mox/token"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, program.String(), "let foo: int = bar;")
}
