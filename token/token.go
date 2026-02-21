package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t *Token) Repr() string {
	return string(t.Type) + "('" + string(t.Literal) + "')"
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "ident"
	INT    = "int"
	FLOAT  = "float"
	STRING = "string"
	BOOL   = "bool"

	DOT = "."

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	FSLASH = "/"
	STAR   = "*"
	LT     = "<"
	GT     = ">"
	BANG   = "!"
	AMP    = "&"
	PIPE   = "|"

	EQ     = "=="
	NOT_EQ = "!="
	GTE    = ">="
	LTE    = "<="
	AND    = "&&"
	OR     = "||"

	PLUS_ASSIGN   = "+="
	MINUS_ASSIGN  = "-="
	STAR_ASSIGN   = "*="
	FSLASH_ASSIGN = "/="

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	QUOTE     = `"`

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LSQB   = "["
	RSQB   = "]"

	// Keywords
	FUNCTION = "fn"
	LET      = "let"
	CLASS    = "class"
	RETURN   = "return"
	FOR      = "for"
	CONST    = "const"
	TRUE     = "true"
	FALSE    = "false"
	IF       = "if"
	IN       = "in"
	ELSE     = "else"
	NULL     = "null"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"class":  CLASS,
	"return": RETURN,
	"for":    FOR,
	"const":  CONST,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"in":     IN,
	"else":   ELSE,
	"null":   NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
