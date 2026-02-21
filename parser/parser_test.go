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
	input := `let x = 122;
	let y: int = 123;
	let z: string = "hello";
	`

	lex := lexer.New(input)
	p := New(lex)
	program := p.ParseProgram()
	stmt := program.Statements
	for i, stmt := range stmt {
		fmt.Printf("line %d ---- %s ---- %+v\n", i, stmt.String(), stmt)
	}

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
				"%d: Expected 'let', foung %q",
				i, stmt.TokenLiteral(),
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

func TestConstStatement(t *testing.T) {
	input := `const x = 122;
	const y: int = 123;
	const z: string = "hello";
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

		if stmt.TokenLiteral() != "const" {
			t.Errorf(
				"%d: Expected 'const', foung %q",
				i, stmt.TokenLiteral(),
			)
		}

		assert.IsType(t, &ast.ConstStatement{}, stmt)
		constStmt, _ := stmt.(*ast.ConstStatement)

		assert.Equal(t, constStmt.Name.Value, tc.expectedIdent)
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
		{
			"add(1, 2)",
			"add(1, 2)",
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

func TestIfStatementStandard(t *testing.T) {
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

func TestIfStatementWithElseIf(t *testing.T) {
	input := `if (x < y) { x } else if (x > y) { z } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	assert.IsType(t, &ast.ExpressionStatement{}, stmt)
	exp := stmt.(*ast.ExpressionStatement).Expression
	is, _ := exp.(*ast.IfExpression)

	assert.Equal(t, "if", is.TokenLiteral())
	assert.Equal(t, "(x < y)", is.Condition.String())

	assert.Equal(t, is.ThenBranch.String(), "x")
	elseIfBranch, _ := is.ElseBranch.(*ast.IfExpression)
	assert.Equal(t, "if", elseIfBranch.TokenLiteral())
	assert.Equal(t, "(x > y)", elseIfBranch.Condition.String())
	assert.Equal(t, "z", elseIfBranch.ThenBranch.String())

	elseBranchStmt, _ := elseIfBranch.ElseBranch.(*ast.BlockStatement).Statements[0].(*ast.ExpressionStatement)
	elseBranchExpr := elseBranchStmt.Expression
	assert.Equal(t, "y", elseBranchExpr.String())
}

func TestIfStatementWithoutBlockStatement(t *testing.T) {
	input := `if (x < y) x else y`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	assert.IsType(t, &ast.ExpressionStatement{}, stmt)
	exp := stmt.(*ast.ExpressionStatement).Expression
	is, _ := exp.(*ast.IfExpression)

	assert.Equal(t, "if", is.TokenLiteral())
	assert.Equal(t, "(x < y)", is.Condition.String())
	assert.Equal(t, "x", is.ThenBranch.String())
	assert.Equal(t, "y", is.ElseBranch.String())
}

func TestIfStatementWithBlockStatementBranch(t *testing.T) {
	input := `if (x < y) { let foo = 3; foo + 1 } else y`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	assert.IsType(t, &ast.ExpressionStatement{}, stmt)
	exp := stmt.(*ast.ExpressionStatement).Expression
	is, _ := exp.(*ast.IfExpression)

	assert.Equal(t, "if", is.TokenLiteral())
	assert.Equal(t, "(x < y)", is.Condition.String())
	assert.Equal(t, "let foo = 3;(foo + 1)", is.ThenBranch.String())
}

func TestFunctionStatementSimple(t *testing.T) {
	input := `fn hello(a, b) { return a+b; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]
	fs, _ := stmt.(*ast.FunctionStatement)

	assert.Equal(t, "fn", fs.TokenLiteral())
	assert.Equal(t, "hello", fs.Name.TokenLiteral())
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
	fs, _ := stmt.(*ast.FunctionStatement)

	assert.Equal(t, "fn", fs.TokenLiteral())
	assert.Equal(t, "helloWithType", fs.Name.TokenLiteral())
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

func TestClassDeclaration(t *testing.T) {
	input := `class Foo(Base) {}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	for _, err := range p.errors {
		fmt.Println(err)
	}

	assert.Equal(t, 1, len(program.Statements))

	stmt := program.Statements[0].(*ast.ClassDeclStatement)
	assert.Equal(t, "class", stmt.TokenLiteral())
	assert.Equal(t, "Foo", stmt.Name.String())
	assert.NotNil(t, stmt.SuperClass)
	assert.Equal(t, "Base", stmt.SuperClass.String())

	assert.Empty(t, stmt.Fields)
	assert.Empty(t, stmt.Methods)
}

func TestClassDeclarationWithFieldsMethods(t *testing.T) {
	input := `class Foo(Base) {
		a: int;
		b: int = 123;
		fn hello(a, b) {
			return a+b;
		}
	}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	for _, err := range p.errors {
		fmt.Println(err)
	}

	assert.Equal(t, 1, len(program.Statements))

	stmt := program.Statements[0].(*ast.ClassDeclStatement)
	assert.Equal(t, "class", stmt.TokenLiteral())
	assert.Equal(t, "Foo", stmt.Name.String())
	assert.NotNil(t, stmt.SuperClass)
	assert.Equal(t, "Base", stmt.SuperClass.String())

	assert.NotEmpty(t, stmt.Fields)
	assert.NotEmpty(t, stmt.Methods)

	assert.Equal(t, 2, len(stmt.Fields))
	assert.Equal(t, "a", stmt.Fields[0].Name.String())
	assert.Equal(t, "int", stmt.Fields[1].Type.String())

	assert.Equal(t, 1, len(stmt.Methods))
	assert.Equal(t, "hello", stmt.Methods[0].Name.String())
	assert.Equal(t, 2, len(stmt.Methods[0].Params))
}

func TestVarStatement(t *testing.T) {
	input := `x: int;`

	l := lexer.New(input)
	p := New(l)

	stmt := p.parseClassVarStatement()
	fmt.Print(p.errors)
	assert.NotNil(t, stmt)

	assert.Equal(t, "x", stmt.Name.String())
	assert.Equal(t, "int", stmt.Type.String())
}

func TestVarStatementWithDefault(t *testing.T) {
	input := `x: int = 123;`

	l := lexer.New(input)
	p := New(l)

	stmt := p.parseClassVarStatement()
	fmt.Print(p.errors)
	assert.NotNil(t, stmt)

	assert.Equal(t, "x", stmt.Name.String())
	assert.Equal(t, "int", stmt.Type.String())
	assert.Equal(t, "123", stmt.Default.String())
}

func TestForInStatementWithExpressionIterable(t *testing.T) {
	input := `for x in getItems() {}`

	l := lexer.New(input)
	p := New(l)

	forIn := p.parseForInStatement()
	assert.IsType(t, &ast.ForInStatement{}, forIn)

	assert.Equal(t, "for", forIn.For.Literal)
	assert.Equal(t, "getItems()", forIn.Iterable.String())
	targets := make([]string, 0)
	for _, t := range forIn.Targets {
		targets = append(targets, t.String())
	}

	assert.ElementsMatch(t, []string{"x"}, targets)
	assert.IsType(t, &ast.BlockStatement{}, forIn.Body)
	assert.Equal(t, 0, len(forIn.Body.Statements))
}

func TestForInStatementMultipleTargets(t *testing.T) {
	input := `for i, x in y {}`

	l := lexer.New(input)
	p := New(l)

	forIn := p.parseForInStatement()
	assert.IsType(t, &ast.ForInStatement{}, forIn)

	assert.Equal(t, "for", forIn.For.Literal)
	assert.Equal(t, "y", forIn.Iterable.String())

	targets := make([]string, 0)
	for _, t := range forIn.Targets {
		targets = append(targets, t.String())
	}

	assert.ElementsMatch(t, []string{"i", "x"}, targets)
	assert.IsType(t, &ast.BlockStatement{}, forIn.Body)
	assert.Equal(t, 0, len(forIn.Body.Statements))
}

func TestForInStatementWithSingleTarget(t *testing.T) {
	input := `for x in y {}`

	l := lexer.New(input)
	p := New(l)

	forIn := p.parseForInStatement()
	assert.IsType(t, &ast.ForInStatement{}, forIn)

	assert.Equal(t, "for", forIn.For.Literal)
	assert.Equal(t, "y", forIn.Iterable.String())
	targets := make([]string, 0)
	for _, t := range forIn.Targets {
		targets = append(targets, t.String())
	}

	assert.ElementsMatch(t, []string{"x"}, targets)
	assert.IsType(t, &ast.BlockStatement{}, forIn.Body)
	assert.Equal(t, 0, len(forIn.Body.Statements))
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, sub(1, 2), foo)`

	program := initParser(input)

	assert.Equal(t, 1, len(program.Statements))
	stmt := program.Statements[0]

	assert.IsType(t, &ast.ExpressionStatement{}, stmt)
	expr, _ := stmt.(*ast.ExpressionStatement)
	c := expr.Expression.(*ast.CallExpression)
	assert.IsType(t, &ast.CallExpression{}, c)

	assert.Equal(t, "(", c.TokenLiteral())
	assert.Equal(t, "add", c.Function.String())
	args := c.Arguments
	assert.Equal(t, "1", args[0].String())
	assert.Equal(t, "(2 * 3)", args[1].String())
	assert.Equal(t, "sub(1, 2)", args[2].String())
	assert.Equal(t, "foo", args[3].String())
}

func TestAssignExpression(t *testing.T) {
	testcases := []struct {
		input    string
		expected string
	}{
		{"x = 5", "x = 5"},
		{"y += 10", "y += 10"},
		{"z -= 3", "z -= 3"},
		{"a *= 2", "a *= 2"},
		{"b /= 4", "b /= 4"},
		{"b /= 4;", "b /= 4"}, // ignores semicolon
	}

	for _, tc := range testcases {
		program := initParser(tc.input)

		assert.Equal(t, 1, len(program.Statements))
		stmt := program.Statements[0]
		assert.IsType(t, &ast.ExpressionStatement{}, stmt)
		expr, _ := stmt.(*ast.ExpressionStatement)
		assignExpr, _ := expr.Expression.(*ast.AssignExpression)

		assert.Equal(t, tc.expected, assignExpr.String())
	}
}

func TestParseMemberAccessExpression(t *testing.T) {
	inputs := []struct {
		input    string
		expected string
	}{
		{"obj.field", "(obj.field)"},
		{"obj.method()", "(obj.method)()"},
		{"obj.field1.field2", "((obj.field1).field2)"},
		{"obj.method1().method2()", "((obj.method1)().method2)()"},
	}

	for _, tc := range inputs {
		program := initParser(tc.input)
		assert.Equal(t, 1, len(program.Statements))
		stmt := program.Statements[0].String()

		assert.Equal(t, tc.expected, stmt)
	}
}

func TestParseIndexExpression(t *testing.T) {
	inputs := []struct {
		input    string
		expected string
	}{
		{"arr[0]", "(arr[0])"},
		{"matrix[1][2]", "((matrix[1])[2])"},
		{"myMap[key]", "(myMap[key])"},
		{"a.b.c[key]", "(((a.b).c)[key])"},
	}

	for _, tc := range inputs {
		program := initParser(tc.input)
		assert.Equal(t, 1, len(program.Statements))
		stmt := program.Statements[0].String()

		assert.Equal(t, tc.expected, stmt)
	}
}
