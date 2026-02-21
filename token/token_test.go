package token

import "testing"

func TestTokenLookupIdent(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"fn", FUNCTION},
		{"let", LET},
		{"class", CLASS},
		{"return", RETURN},
		{"for", FOR},
		{"const", CONST},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"in", IN},
		{"else", ELSE},
		{"null", NULL},
		{"map", MAP},
		{"array", ARRAY},
		{"set", SET},
		{"tuple", TUPLE},
		{"union", UNION},
		{"option", OPTION},

		{"foobar", IDENT},
	}

	for _, tt := range tests {
		if tok := LookupIdent(tt.input); tok != tt.expected {
			t.Errorf("LookupIdent(%q) = %q, expected %q", tt.input, tok, tt.expected)
		}
	}
}
