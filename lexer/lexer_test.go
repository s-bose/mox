package lexer

import (
	"testing"

	"github.com/s-bose/mox/token"
	"github.com/stretchr/testify/assert"
)

type TokenExpectation struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	const x: string = "hello";`

	tests := []TokenExpectation{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.CONST, "const"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.IDENT, "string"},
		{token.ASSIGN, "="},
		{token.STRING, "hello"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	for _, tt := range tests {
		tok := l.NextToken()

		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestNextTokenOperator(t *testing.T) {
	input := `2 < 3
	<=
	>=
	!=
	==
	&&
	||
	`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.INT, "2"},
		{token.LT, "<"},
		{token.INT, "3"},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.NOT_EQ, "!="},
		{token.EQ, "=="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tok.Literal, tt.expectedLiteral)
		assert.Equal(t, tok.Type, tt.expectedType)
	}
}
