package parser

import (
	"fmt"
	"testing"

	"github.com/s-bose/mox/ast"
	"github.com/s-bose/mox/lexer"
	"github.com/stretchr/testify/assert"
)

func initParser(input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
	}

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestLetStatements(t *testing.T) {
	input := `
	let x = 123;
	let y: int = 123;
	let z: string = "hello";
	`

	lex := lexer.New(input)
	p := New(lex)
	program := p.ParseProgram()

	tests := []struct {
		expectedIdent string
		expectedType  string
	}{
		{"x", ""},
		{"y", "int"},
		{"z", "string"},
	}

	assert.NotNil(t, program)
	assert.Equal(t, len(program.Statements), 3)
	for i, tc := range tests {
		stmt := program.Statements[i]

		if stmt.TokenLiteral() != "let" {
			t.Errorf(
				"Expected 'let', foung %q",
				stmt.TokenLiteral(),
			)
		}

		assert.IsType(t, &ast.LetStatement{}, stmt)
		letStmt, _ := stmt.(*ast.LetStatement)

		assert.Equal(t, letStmt.Name.Value, tc.expectedIdent)

		if letStmt.Type != nil {
			assert.Equal(t, letStmt.Type.Value, tc.expectedType)
		}

		assert.Equal(t, letStmt.Name.TokenLiteral(), tc.expectedIdent)
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

	assert.Equal(t, len(program.Statements), 3)
	for _, stmt := range program.Statements {
		assert.IsType(t, &ast.ReturnStatement{}, stmt)

		returnStmtStruct, _ := stmt.(*ast.ReturnStatement)

		assert.Equal(t, returnStmtStruct.TokenLiteral(), "return")
	}
}

func TestIntegerExpression(t *testing.T) {
	inputs := []struct {
		inputStr    string
		expectedInt int64
	}{
		{"1;", 1},
		{"1234;", 1234},
	}

	for _, tt := range inputs {
		l := lexer.New(tt.inputStr)
		p := New(l)

		program := p.ParseProgram()
		assert.Equal(t, len(program.Statements), 1)
		stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
		intExpr, _ := stmt.Expression.(*ast.IntegerLiteral)
		assert.Equal(t, intExpr.Value, tt.expectedInt)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foo;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	assert.Equal(t, len(program.Statements), 1)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
	ident, _ := stmt.Expression.(*ast.Identifier)

	assert.Equal(t, ident.Value, "foo")
	assert.Equal(t, ident.TokenLiteral(), "foo")
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
		assert.Equal(t, len(program.Statements), 1)
		stmt, _ := program.Statements[0].(*ast.ExpressionStatement)
		exp, _ := stmt.Expression.(*ast.PrefixExpr)

		assert.Equal(t, exp.Op, tt.op)
		rval, _ := exp.Right.(*ast.IntegerLiteral)
		assert.Equal(t, rval.Value, tt.intValue)
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
		{
			"3 + 4 <= 5",
			"((3 + 4) <= 5)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"2 < 3 == true",
			"((2 < 3) == true)",
		},
		{
			"3 < 2 == false",
			"((3 < 2) == false)",
		},
		{
			"1 + (2 + 3)",
			"(1 + (2 + 3))",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)

		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()

		assert.Equal(t, actual, tt.expected)
	}
}

func TestIfStatement(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	assert.Equal(t, len(program.Statements), 1)
	stmt := program.Statements[0]
	assert.IsType(t, &ast.ExpressionStatement{}, stmt)

	exp := stmt.(*ast.ExpressionStatement).Expression
	is, _ := exp.(*ast.IfExpression)
	assert.Equal(t, "if", is.TokenLiteral())
	assert.Equal(t, "(x < y)", is.Condition.String())
	assert.Equal(t, "x", is.ThenBranch.String())
	assert.Equal(t, "y", is.ElseBranch.String())
}

func TestIfStatementWithElseIfStatement(t *testing.T) {
	input := `if (x <= y) { x } else if (x < y) { z } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	assert.IsType(t, &ast.ExpressionStatement{}, stmt)
	exp := stmt.(*ast.ExpressionStatement).Expression
	is, _ := exp.(*ast.IfExpression)

	assert.Equal(t, "if", is.TokenLiteral())
	assert.Equal(t, "(x <= y)", is.Condition.String())

	assert.Equal(t, is.ThenBranch.String(), "x")
	elseBranchStmt, _ := is.ElseBranch.Statements[0].(*ast.ExpressionStatement)
	elseExpr := elseBranchStmt.Expression
	assert.IsType(t, &ast.IfExpression{}, elseExpr)

	elseIfExpr, _ := elseExpr.(*ast.IfExpression)

	assert.Equal(t, "if", elseIfExpr.TokenLiteral())
	assert.Equal(t, "(x < y)", elseIfExpr.Condition.String())
	assert.Equal(t, "z", elseIfExpr.ThenBranch.String())
}

func TestFunctionExpressionSimple(t *testing.T) {
	input := `fn hello(a, b) { return a+b; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	exp := stmt.(*ast.ExpressionStatement).Expression
	fs, _ := exp.(*ast.FunctionLiteral)

	assert.Equal(t, "hello", fs.TokenLiteral())
	fmt.Println(fs)
	assert.Equal(t, 2, len(fs.Params))

	assert.Equal(t, "a", fs.Params[0].Token.Literal)
	assert.Equal(t, "b", fs.Params[1].Token.Literal)

	assert.Empty(t, fs.Defaults)
	assert.Empty(t, fs.ParamType)
	assert.Nil(t, fs.ReturnType)

	assert.NotNil(t, fs.Body)
	body := fs.Body
	assert.Equal(t, 1, len(body.Statements))
	bodyStmt := body.Statements[0]

	returnStmt := bodyStmt.(*ast.ReturnStatement)
	assert.Equal(t, "return", returnStmt.TokenLiteral())
	assert.Equal(t, "(a + b)", returnStmt.ReturnValue.String())
}

func TestFunctionExpressionWithTypesAndDefault(t *testing.T) {
	input := `fn helloWithType(a: int, b: string, c: int = 10): string { return b; }`

	program := initParser(input)

	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	exp := stmt.(*ast.ExpressionStatement).Expression
	fs, _ := exp.(*ast.FunctionLiteral)

	assert.Equal(t, "helloWithType", fs.TokenLiteral())
	fmt.Println(fs)
	assert.Equal(t, 3, len(fs.Params))

	assert.Equal(t, "a", fs.Params[0].Token.Literal)
	assert.Equal(t, "b", fs.Params[1].Token.Literal)
	assert.Equal(t, "c", fs.Params[2].Token.Literal)

	assert.NotEmpty(t, fs.Defaults)
	assert.Equal(t, "10", fs.Defaults["c"].TokenLiteral())
	assert.Nil(t, fs.Defaults["a"])
	assert.Nil(t, fs.Defaults["b"])

	assert.NotEmpty(t, fs.ParamType)
	assert.Equal(t, "int", fs.ParamType["a"].TokenLiteral())
	assert.Equal(t, "string", fs.ParamType["b"].TokenLiteral())
	assert.Equal(t, "int", fs.ParamType["c"].TokenLiteral())

	assert.NotNil(t, fs.ReturnType)
	assert.Equal(t, "string", fs.ReturnType.TokenLiteral())

	assert.NotNil(t, fs.Body)
	body := fs.Body
	assert.Equal(t, 1, len(body.Statements))
	bodyStmt := body.Statements[0]

	returnStmt := bodyStmt.(*ast.ReturnStatement)
	assert.Equal(t, "return", returnStmt.TokenLiteral())
	assert.Equal(t, "b", returnStmt.ReturnValue.String())
}
