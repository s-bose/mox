package parser

import (
	"fmt"
	"strconv"

	"github.com/s-bose/mox/ast"
	"github.com/s-bose/mox/lexer"
	"github.com/s-bose/mox/token"
)

type Precedence int

const (
	_           Precedence = iota
	LOWEST                 // _
	EQUALS                 // ==
	LESSGREATER            // < >
	SUM                    // +
	PRODUCT                // *
	PREFIX                 // -X or !X
	CALL                   // abc(X)
)

var precedenceTable = map[token.TokenType]Precedence{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.LTE:    LESSGREATER,
	token.GTE:    LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.STAR:   PRODUCT,
	token.FSLASH: PRODUCT,
	token.LPAREN: CALL,
}

type (
	TPrefixParseFunc func() ast.Expression
	TInfixParseFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFuncs map[token.TokenType]TPrefixParseFunc
	infixParseFuncs  map[token.TokenType]TInfixParseFunc
}

func initParseFuncs(p *Parser) {
	// Prefix fns
	p.registerPrefixFunc(token.IDENT, p.parseIdentifier)
	p.registerPrefixFunc(token.INT, p.parseIntegerExpr)
	p.registerPrefixFunc(token.FLOAT, p.parseFloatExpr)
	p.registerPrefixFunc(token.TRUE, p.parseBooleanExpr)
	p.registerPrefixFunc(token.FALSE, p.parseBooleanExpr)
	p.registerPrefixFunc(token.LPAREN, p.parseGroupedExpr)
	p.registerPrefixFunc(token.IF, p.parseIfExpr)

	p.registerPrefixFunc(token.BANG, p.parsePrefixExpr)
	p.registerPrefixFunc(token.MINUS, p.parsePrefixExpr)

	// Infix fns
	p.registerInfixFunc(token.PLUS, p.parseInfixExpr)
	p.registerInfixFunc(token.MINUS, p.parseInfixExpr)
	p.registerInfixFunc(token.FSLASH, p.parseInfixExpr)
	p.registerInfixFunc(token.STAR, p.parseInfixExpr)
	p.registerInfixFunc(token.EQ, p.parseInfixExpr)
	p.registerInfixFunc(token.NOT_EQ, p.parseInfixExpr)
	p.registerInfixFunc(token.LT, p.parseInfixExpr)
	p.registerInfixFunc(token.GT, p.parseInfixExpr)
	p.registerInfixFunc(token.LTE, p.parseInfixExpr)
	p.registerInfixFunc(token.GTE, p.parseInfixExpr)
	p.registerInfixFunc(token.LPAREN, p.parseCallExpr)
}
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFuncs = make(map[token.TokenType]TPrefixParseFunc)
	p.infixParseFuncs = make(map[token.TokenType]TInfixParseFunc)

	initParseFuncs(p)
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekError(t token.TokenType) {
	s := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, s)
}

func (p *Parser) peekPrecedence() Precedence {
	if pre, ok := precedenceTable[p.peekToken.Type]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) curPrecedence() Precedence {
	if pre, ok := precedenceTable[p.curToken.Type]; ok {
		return pre
	}
	return LOWEST
}

func (p *Parser) addParseError(s string) {
	p.errors = append(p.errors, s)
}

// register parse fns
func (p *Parser) registerPrefixFunc(token token.TokenType, fn TPrefixParseFunc) {
	p.prefixParseFuncs[token] = fn
}

func (p *Parser) registerInfixFunc(token token.TokenType, fn TInfixParseFunc) {
	p.infixParseFuncs[token] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// Expression Parsers

func (p *Parser) parseBooleanExpr() ast.Expression {
	b := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}

	return b
}

func (p *Parser) parseIntegerExpr() ast.Expression {
	il := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addParseError(fmt.Sprintf("could not parse %s as integer: %v", p.curToken.Literal, err))
	}
	il.Value = value
	return il
}

func (p *Parser) parseFloatExpr() ast.Expression {
	fl := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addParseError(fmt.Sprintf("could not parse %s as float: %v", p.curToken.Literal, err))
	}

	fl.Value = value
	return fl
}

func (p *Parser) parseGroupedExpr() ast.Expression {
	if p.peekTokenIs(token.RPAREN) {
		// nothing to parse, empty paren
		return nil
	}

	p.nextToken()
	exp := p.parseExpr(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		statement := p.ParseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parsePrefixExpr() ast.Expression {
	expr := &ast.PrefixExpr{
		Token: p.curToken,
		Op:    p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpr(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpr{
		Token: p.curToken,
		Op:    p.curToken.Literal,
		Left:  left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpr(precedence)

	return expr
}

func (p *Parser) parseExpr(precedence Precedence) ast.Expression {
	prefixFunc := p.prefixParseFuncs[p.curToken.Type]
	if prefixFunc == nil {
		p.addParseError(fmt.Sprintf("no prefix function exists for %s", p.curToken.Type))
		return nil
	}

	leftExp := prefixFunc()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFunc := p.infixParseFuncs[p.peekToken.Type]
		if infixFunc == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infixFunc(leftExp)
	}
	return leftExp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := make([]ast.Expression, 0)
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpr(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpr(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		p.addParseError(fmt.Sprintf("expected closing ')' for arguments, got %s", p.peekToken.Literal))
		return nil
	}

	return args
}

func (p *Parser) parseCallExpr(function ast.Expression) ast.Expression {
	stmt := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	stmt.Arguments = p.parseCallArguments()
	return stmt
}

// Parse statements helpers

func (p *Parser) parseForInStatement() *ast.ForInStatement {
	stmt := &ast.ForInStatement{
		For: p.curToken,
	}
	p.nextToken()

	targets := make([]*ast.Identifier, 0)

	for !p.curTokenIs(token.IN) {
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		} else if p.curTokenIs(token.IDENT) {
			targetIdent := &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			}

			targets = append(targets, targetIdent)
			p.nextToken()
		} else {
			p.addParseError(fmt.Sprintf("invalid token %s", p.curToken.Literal))
			return nil
		}
	}

	if len(targets) != 0 {
		stmt.Targets = targets
	}

	// consume `in`
	p.nextToken()

	stmt.Iterable = p.parseExpr(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		p.addParseError("for loop missing block statement")
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	// parse for statement
	// example
	//
	// for (let x = 1; x < 10; x += 1) { ... }
	//
	// for i, iter in iterable { ... }

	// forIdent := &ast.Identifier{
	// 	Token: p.curToken,
	// 	Value: p.curToken.Literal,
	// }

	if p.peekTokenIs(token.LPAREN) {
		// p.parseForStatementWithInitCond
	}

	if p.peekTokenIs(token.IDENT) {
		return p.parseForInStatement()
	}

	return nil
}

func (p *Parser) parseClassVarStatement() *ast.ClassVarStatement {
	/* Parses the following rule
	 * <identifier>: <type>? = <default>?
	 * This is only an allowed syntax inside thr class definition
	 */
	if !p.curTokenIs(token.IDENT) {
		p.addParseError(fmt.Sprintf("expected identifier, got %s", p.curToken.Literal))
		return nil
	}

	stmt := &ast.ClassVarStatement{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	if !p.expectPeek(token.COLON) {
		p.addParseError(fmt.Sprintf("expected `:` after identifier, got %s", p.peekToken.Literal))
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		p.addParseError(fmt.Sprintf("expected type specifier after `:`, got %s", p.curToken.Literal))
		return nil
	}

	stmt.Type = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	p.nextToken()

	if p.curTokenIs(token.ASSIGN) {
		// optional default assignment
		p.nextToken()
		stmt.Default = &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		p.nextToken()
	}

	if !p.curTokenIs(token.SEMICOLON) {
		p.addParseError(fmt.Sprintf("var definitions should end in semicolons, got %s", p.curToken.Literal))
		return nil
	}

	return stmt
}

func (p *Parser) parseClassBlockStatements() ([]*ast.ClassVarStatement, []*ast.FunctionStatement) {
	fields := make([]*ast.ClassVarStatement, 0)
	methods := make([]*ast.FunctionStatement, 0)

	p.nextToken()

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) {
			field := p.parseClassVarStatement()
			if field != nil {
				fields = append(fields, field)
			}
			fmt.Print(fields, "\n", p.curToken.Type)
		} else if p.curTokenIs(token.FUNCTION) {
			meth := p.parseFunctionStatement()
			if meth != nil {
				methods = append(methods, meth)
			}
			fmt.Print(fields)
		} else {
			p.addParseError(fmt.Sprintf("invalid keyword %s inside class definition", p.curToken.Literal))
			return fields, methods
		}

		p.nextToken()
	}

	return fields, methods
}

func (p *Parser) parseClassDeclaration() *ast.ClassDeclStatement {
	stmt := &ast.ClassDeclStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		p.addParseError(fmt.Sprintf("expected class name, found %s", p.peekToken.Literal))
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			stmt.SuperClass = &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			}
		}

		if !p.expectPeek(token.RPAREN) {
			p.addParseError(fmt.Sprintf("expected ')', found %s", p.peekToken.Literal))
			return nil
		}
	}

	// }
	if !p.expectPeek(token.LBRACE) {
		p.addParseError("class definition missing '{'")
		return nil
	}

	fmt.Printf("curToken: %s, nextToken: %s\n", p.curToken.Literal, p.peekToken.Literal)
	stmt.Fields, stmt.Methods = p.parseClassBlockStatements()
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	/*
	* Parses a return statement
	*
	* Example
	* return x;
	* return x / y;
	* return 12;
	* return <expr>;
	 */

	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpr(LOWEST)

<<<<<<< Updated upstream
	// skip next tokens until we hit semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	/*
		 	* Parses a let statement
			*
			* Example
			* let x = 123;
			* let y: int = 42;
			* let z: string = "hello";
	*/

	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.expectPeek(token.COLON) {
		// parse variable type
		p.nextToken()

		typeIdent := token.LookupDataType(p.curToken.Literal)
		p.curToken.Type = typeIdent
		stmt.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

<<<<<<< Updated upstream
	stmt.Value = p.parseExpr(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		// skip tokens until semicolon (for now)
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	/*
		 *
			* Parses an expression statement
			*
			*
			* Examples
			* ## Prefix expressions
			*
			* -4 -variable_name, !func()
			*
			* ## Infix expressions
			*
			* x + - * / y
			* x == y
			* x <= y
			* x >= y
			* x < y
			* x > y
			* x != y
	*/
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpr(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.curToken,
	}

	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionParameters() ([]*ast.Identifier, map[string]*ast.Identifier, map[string]ast.Expression) {
	argumentTypes := make(map[string]*ast.Identifier)
	argumentDefaults := make(map[string]ast.Expression)

	idents := make([]*ast.Identifier, 0)
	if p.peekTokenIs(token.RPAREN) {
		// return parsed params on closing braces: ")"
		p.nextToken()
		return idents, argumentTypes, argumentDefaults
	}

	p.nextToken()
	for !p.curTokenIs(token.RPAREN) {
		if p.curTokenIs(token.EOF) {
			p.addParseError("unexpected EOF on function parameters")
			return nil, nil, nil
		}

		ident := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		idents = append(idents, ident)
		p.nextToken()

		if p.curTokenIs(token.COLON) {
			// ident: type
			p.nextToken() // get the type token after `:`
			typeIdent := &ast.Identifier{
				Token: p.curToken,
				Value: p.curToken.Literal,
			}
			argumentTypes[ident.Value] = typeIdent
			p.nextToken()
		}

		if p.curTokenIs(token.ASSIGN) {
			p.nextToken()
			argumentDefaults[ident.Value] = p.parseExpressionStatement().Expression
			p.nextToken()
		}

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return idents, argumentTypes, argumentDefaults
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	/*
		 * Parses a function statement
			* functionStatement :== fn <Identifier> "(" <args...>, <args: type>, <args: type "=" <Identifier>> ): <returnTypeIdentifier>
			* 						"{" <blockStatement> "}"
		 	* Examples
				* fn foobar(a: int, b: string, c: MyClass = DEFAULT) {
				* 	return a + b;
				* }
	*/

	// reference example: fn hello(a: int = 100): string { ... }

	// -- fn
	stmt := &ast.FunctionStatement{
		Token: p.curToken, // stores `hello`
	}
	// -- hello
	if !p.expectPeek(token.IDENT) {
		p.addParseError("expected function name, got none")
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	// -- (
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// [ a ], {a: int}, {a: 100}
	args, argTypes, argDefaults := p.parseFunctionParameters()
	stmt.Params = args
	stmt.ParamType = argTypes
	stmt.Defaults = argDefaults

	// -- )
	if !p.curTokenIs(token.RPAREN) {
		return nil
	}

	// -- : string
	if p.peekTokenIs(token.COLON) {
		// variable : type
		p.nextToken() // consume `:`
		p.nextToken() // type(string)
		retTypeIdent := &ast.Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		stmt.ReturnType = retTypeIdent
	}

	// {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	p.nextToken()
	return stmt
}

func (p *Parser) parseBracketExpression() ast.Expression {
	if !p.expectPeek(token.LPAREN) {
		p.addParseError(fmt.Sprintf("expected '(' but got %s", p.curToken.Literal))
		return nil
	}

	p.nextToken()
	expr := p.parseExpr(LOWEST)
	if expr == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		p.addParseError(fmt.Sprintf("expected ')' but got %s", p.curToken.Literal))
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpr() ast.Expression {
	stmt := &ast.IfExpression{
		Token: p.curToken,
	}

	cond := p.parseBracketExpression()
	if cond == nil {
		return nil
	}
	stmt.Condition = cond

	fmt.Println(stmt.Condition.String())

	// parse then clause { ... }

	if !p.expectPeek(token.LBRACE) {
		p.addParseError(fmt.Sprintf("expected '{' but got %s", p.curToken.Literal))
		return nil
	}

	blockStmt := p.parseBlockStatement()
	stmt.ThenBranch = blockStmt

	// else block
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		// `else if`
		if p.peekTokenIs(token.IF) {
			p.nextToken()

			stmt.ElseBranch = &ast.BlockStatement{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Expression: p.parseIfExpr(),
					},
				},
			}

			return stmt
		}

		if !p.expectPeek(token.LBRACE) {
			p.addParseError(fmt.Sprintf("expected '{' but got %s", p.curToken.Literal))
			return nil
		}

		stmt.ElseBranch = p.parseBlockStatement()
	}
	return stmt
}

func (p *Parser) ParseStatement() ast.Statement {
	fmt.Printf("curToken is: %s\n", p.curToken.Literal)

	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.CLASS:
		return p.parseClassDeclaration()
	case token.FOR:
		return p.parseForStatement()
	default:
		return p.parseExpressionStatement()
	}
}
