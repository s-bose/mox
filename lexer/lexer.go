package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/s-bose/mox/token"
)

type Lexer struct {
	position     int
	readPosition int
	ch           rune
	input        []rune
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition

	l.readPosition += 1
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case rune('#'):
		// skip single-line comments starting with #
		l.skipComment()
		return l.NextToken()
	case rune('='):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = makeToken(token.ASSIGN, l.ch)
		}
	case rune('!'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = makeToken(token.BANG, l.ch)
		}
	case rune('+'):
		tok = makeToken(token.PLUS, l.ch)
	case rune('-'):
		tok = makeToken(token.MINUS, l.ch)
	case rune('/'):
		tok = makeToken(token.FSLASH, l.ch)
	case rune('*'):
		tok = makeToken(token.STAR, l.ch)
	case rune('<'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = makeToken(token.LT, l.ch)
		}
	case rune('>'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = makeToken(token.GT, l.ch)
		}
	case rune('&'):
		if l.peekChar() == rune('&') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			tok = makeToken(token.AMP, l.ch)
		}
	case rune('|'):
		if l.peekChar() == rune('|') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			tok = makeToken(token.PIPE, l.ch)
		}
	case rune(','):
		tok = makeToken(token.COMMA, l.ch)
	case rune(';'):
		tok = makeToken(token.SEMICOLON, l.ch)
	case rune(':'):
		tok = makeToken(token.COLON, l.ch)
	case rune('"'):
		str, err := l.readString()
		if err != nil {
			tok.Literal = err.Error()
			tok.Type = token.ILLEGAL
		}
		tok.Literal = str
		tok.Type = token.STRING
	case rune('('):
		tok = makeToken(token.LPAREN, l.ch)
	case rune(')'):
		tok = makeToken(token.RPAREN, l.ch)
	case rune('{'):
		tok = makeToken(token.LBRACE, l.ch)
	case rune('}'):
		tok = makeToken(token.RBRACE, l.ch)
	case rune('['):
		tok = makeToken(token.LSQB, l.ch)
	case rune(']'):
		tok = makeToken(token.RSQB, l.ch)
	case rune('.'):
		tok = makeToken(token.DOT, l.ch)
	case rune(0):
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:

		if isDigit(l.ch) {
			tok := l.readDecimal()
			return tok
		}

		if !unicode.IsPrint(l.ch) {
			tok = makeToken(token.ILLEGAL, l.ch)
			l.readChar()
			return tok
		}

		tok.Literal = l.readIdent()

		if len(tok.Literal) == 0 {
			tok.Type = token.ILLEGAL
			l.readChar()
			return tok
		}

		tok.Type = token.LookupIdent(tok.Literal)
		return tok

	}

	l.readChar()
	return tok
}

func (l *Lexer) skipComment() {
	for l.ch != rune('\n') && l.ch != rune(0) {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) readIdent() string {
	// first character must be a letter or underscore
	if !isIdentStart(l.ch) {
		return ""
	}

	id := string(l.ch)
	l.readChar()

	for isIdentPart(l.ch) {
		id += string(l.ch)
		l.readChar()
	}

	return id
}

func (l *Lexer) readString() (string, error) {
	out := ""

	for {
		l.readChar()

		if l.ch == rune(0) {
			return "", fmt.Errorf("string not terminated")
		}

		if l.ch == '"' {
			break
		}

		if l.ch == '\\' {
			if l.peekChar() == '\n' {
				l.readChar()
				continue
			}

			l.readChar()
			if l.ch == 0 {
				return "", fmt.Errorf("string not terminated")
			}
			if l.ch == 'n' {
				l.ch = '\n'
			}
			if l.ch == 'r' {
				l.ch = '\r'
			}
			if l.ch == 't' {
				l.ch = '\t'
			}
			if l.ch == '"' {
				l.ch = '"'
			}
			if l.ch == '\\' {
				l.ch = '\\'
			}
		}
		out += string(l.ch)
	}
	return out, nil
}

func (l *Lexer) readNumber() string {
	str := ""

	accept := "0123456789"

	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		accept = "0x123456789abcdefABCDEF"
	}

	if l.ch == '0' && (l.peekChar() == 'b' || l.peekChar() == 'B') {
		accept = "b01"
	}

	for strings.Contains(accept, string(l.ch)) {
		str += string(l.ch)
		l.readChar()
	}

	return str
}

func (l *Lexer) readDecimal() token.Token {
	integer := l.readNumber()

	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		fraction := l.readNumber()
		return token.Token{Type: token.FLOAT, Literal: integer + "." + fraction}
	}
	return token.Token{Type: token.INT, Literal: integer}
}

func (l *Lexer) skipWhitespace() {

	for l.ch == rune(' ') || l.ch == rune('\t') || l.ch == rune('\r') || l.ch == rune('\n') {
		l.readChar()
	}
}

func isDigit(ch rune) bool {
	return rune('0') <= ch && ch <= rune('9')
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_' || unicode.IsDigit(ch)
}

func makeToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
