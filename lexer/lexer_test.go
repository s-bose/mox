package lexer

import (
	"testing"

	"github.com/s-bose/mox/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	const x: string = "hello world";`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.CONST, "const"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.STRING, "string"},
		{token.ASSIGN, "="},
		{token.QUOTE, `"`},
		{token.IDENT, "hello world"},
		{token.QUOTE, `"`},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("tests[%d] - token type is not %q. got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Errorf("tests[%d] - token literal is not %q. got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
