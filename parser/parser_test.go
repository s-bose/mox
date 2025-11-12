package parser

import (
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	input := "foo;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, expected %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program missing expected statement, %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("expression not *ast.Identifier, got %T", stmt.Expression)
	}
	if ident.Value != "foo" {
		t.Fatalf("ident.Value not %s, got %s", "foo", ident.Value)
	}

	if ident.TokenLiteral() != "foo" {
		t.Fatalf("ident.TokenLiteral() not %s, got %s", "foo", ident.TokenLiteral())
	}

}

func TestPrefixExpression(t *testing.T) {
	testCases := []struct {
		input    string
		op       string
		intValue int64
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
	}

	for _, tt := range testCases {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		if len(program.Statements) != 1 {
			t.Errorf("program.Statements contains %d statements, expected 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is of type %T, expected ast.ExpressionStatement", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt.Expression is of type %T, expected ast.PrefixExpr", stmt.Expression)
		}

		if exp.Op != tt.op {
			t.Fatalf("expected op to be %s, got %s", tt.op, exp.Op)
		}

		rval, _ := exp.Right.(*ast.IntegerLiteral)
		if rval.Value != tt.intValue {
			t.Fatalf("expected integer value to be %d, got %d", tt.intValue, rval.Value)
		}
	}
}

func TestOpPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"-a-b",
			"((-a) - b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		p := New(l)
		program := p.ParseProgram()

		fmt.Print(program.Statements[0])
		actual := program.String()

		if actual != tt.expected {
			t.Fatalf("expected statement to be %s, got %s", tt.expected, actual)
		}
	}
}
