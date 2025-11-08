package parser

import (
	"testing"

	"github.com/s-bose/mox/ast"
	"github.com/s-bose/mox/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 123;
	let y: int = 123;
	let z: string = "hello";
	`

	lex := lexer.New(input)
	p := New(lex)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("program does not contain any statements, expected %d", len(program.Statements))
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program should contain %d statements", len(program.Statements))
	}

	tests := []struct {
		expectedIdent string
		expectedType  string
	}{
		{"x", ""},
		{"y", "int"},
		{"z", "string"},
	}

	for i, tc := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != "let" {
			t.Errorf(
				"Expected 'let', foung %q",
				stmt.TokenLiteral(),
			)
		}

		letStmt, ok := stmt.(*ast.LetStatement)
		if !ok {
			t.Errorf("s is not of type LetStatement, %T", letStmt)
		}

		if letStmt.Name.Value != tc.expectedIdent {
			t.Fatalf("expected identifier %s, found %s", tc.expectedIdent, letStmt.Name.Value)
		}

		if letStmt.Type != nil {
			if letStmt.Type.Value != tc.expectedType {
				t.Fatalf("expected type %s, found %s", tc.expectedType, letStmt.Type.Value)
			}
		}

		if letStmt.Name.TokenLiteral() != tc.expectedIdent {
			t.Fatalf("")
		}

	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 1;
	return 2;
	return 1234;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmtStruct, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Fatalf("unable to parse return statement, %T", stmt)
			continue
		}

		if returnStmtStruct.TokenLiteral() != "return" {
			t.Fatalf("invalid token, got %q", returnStmtStruct.TokenLiteral())
		}
	}
}
