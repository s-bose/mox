package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "int"
	FLOAT  = "float"
	STRING = "string"
	BOOL   = "bool"

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	FSLASH = "/"
	STAR   = "*"

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
	ELSE     = "else"
	NULL     = "null"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"class":  CLASS,
	"return": RETURN,
	"for":    FOR,
	"int":    INT,
	"float":  FLOAT,
	"string": STRING,
	"bool":   BOOL,
	"const":  CONST,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"null":   NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
