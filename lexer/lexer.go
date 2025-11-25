package lexer

import (
	"fmt"

	"github.com/s-bose/mox/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
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

func (l *Lexer) PeekChar() byte {
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
	case '=':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = makeToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = makeToken(token.BANG, l.ch)
		}
	case '+':
		tok = makeToken(token.PLUS, l.ch)
	case '-':
		tok = makeToken(token.MINUS, l.ch)
	case '/':
		tok = makeToken(token.FSLASH, l.ch)
	case '*':
		tok = makeToken(token.STAR, l.ch)
	case '<':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = makeToken(token.LT, l.ch)
		}
	case '>':
		if l.PeekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = makeToken(token.GT, l.ch)
		}
	case '&':
		if l.PeekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			tok = makeToken(token.AMP, l.ch)
		}
	case '|':
		if l.PeekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			tok = makeToken(token.PIPE, l.ch)
		}
	case ',':
		tok = makeToken(token.COMMA, l.ch)
	case ';':
		tok = makeToken(token.SEMICOLON, l.ch)
	case ':':
		tok = makeToken(token.COLON, l.ch)
	case '"':
		str, err := l.readString()
		if err != nil {
			tok.Literal = err.Error()
			tok.Type = token.ILLEGAL
		}
		tok.Literal = str
		tok.Type = token.STRING
	case '(':
		tok = makeToken(token.LPAREN, l.ch)
	case ')':
		tok = makeToken(token.RPAREN, l.ch)
	case '{':
		tok = makeToken(token.LBRACE, l.ch)
	case '}':
		tok = makeToken(token.RBRACE, l.ch)
	case '[':
		tok = makeToken(token.LSQB, l.ch)
	case ']':
		tok = makeToken(token.RSQB, l.ch)
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = makeToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdent() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, error) {
	out := ""

	for {
		l.readChar()

		if l.ch == 0 {
			return "", fmt.Errorf("string not terminated")
		}

		if l.ch == '"' {
			break
		}

		if l.ch == '\\' {
			if l.PeekChar() == '\n' {
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
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func makeToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
