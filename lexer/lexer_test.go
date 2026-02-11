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
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestDotToken(t *testing.T) {
	input := `foo.bar.baz`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.IDENT, "foo"},
		{token.DOT, "."},
		{token.IDENT, "bar"},
		{token.DOT, "."},
		{token.IDENT, "baz"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestDotTokenWithMethodCall(t *testing.T) {
	input := `obj.method(x, y)`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.IDENT, "obj"},
		{token.DOT, "."},
		{token.IDENT, "method"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestFloatLiteral(t *testing.T) {
	input := `3.14`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.FLOAT, "3.14"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestIntegerFormats(t *testing.T) {
	input := `42 0xFF 0b1010`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.INT, "42"},
		{token.INT, "0xFF"},
		{token.INT, "0b1010"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestFloatVsDotAccess(t *testing.T) {
	input := `a.b 1.5 c.d`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.IDENT, "a"},
		{token.DOT, "."},
		{token.IDENT, "b"},
		{token.FLOAT, "1.5"},
		{token.IDENT, "c"},
		{token.DOT, "."},
		{token.IDENT, "d"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestStringEscapes(t *testing.T) {
	input := `"hello\nworld" "tab\there"`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.STRING, "hello\nworld"},
		{token.STRING, "tab\there"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestKeywordsAndIdentifiers(t *testing.T) {
	input := `fn myFunc let x class Foo return if else for in const true false null`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.FUNCTION, "fn"},
		{token.IDENT, "myFunc"},
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.CLASS, "class"},
		{token.IDENT, "Foo"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.FOR, "for"},
		{token.IN, "in"},
		{token.CONST, "const"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NULL, "null"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestDataTypes(t *testing.T) {
	input := `let x: int = 5; let y: float = 3.14; let z: bool = true;`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.IDENT, "int"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.COLON, ":"},
		{token.IDENT, "float"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.14"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "z"},
		{token.COLON, ":"},
		{token.IDENT, "bool"},
		{token.ASSIGN, "="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}

func TestBracketsAndBraces(t *testing.T) {
	input := `[1, 2] {a: b} (x)`

	l := New(input)

	expectedTokens := []TokenExpectation{
		{token.LSQB, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RSQB, "]"},
		{token.LBRACE, "{"},
		{token.IDENT, "a"},
		{token.COLON, ":"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	for _, tt := range expectedTokens {
		tok := l.NextToken()
		assert.Equal(t, tt.expectedType, tok.Type)
		assert.Equal(t, tt.expectedLiteral, tok.Literal)
	}
}
