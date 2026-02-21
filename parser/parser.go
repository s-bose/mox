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
	_      Precedence = iota
	LOWEST            // _
	ASSIGN            // = += -= *= /=
	EQUALS            // ==

	LESSGREATER  // < or >
	SUM          // + or -
	PRODUCT      // / or *
	PREFIX       // -X or !X
	MEMBERACCESS // foo.bar.baz
	INDEX        // array[index]
	CALL         // abc(X)
)

var precedenceTable = map[token.TokenType]Precedence{
	token.EQ:            EQUALS,
	token.NOT_EQ:        EQUALS,
	token.ASSIGN:        ASSIGN,
	token.PLUS_ASSIGN:   ASSIGN,
	token.MINUS_ASSIGN:  ASSIGN,
	token.STAR_ASSIGN:   ASSIGN,
	token.FSLASH_ASSIGN: ASSIGN,
	token.LT:            LESSGREATER,
	token.GT:            LESSGREATER,
	token.LTE:           LESSGREATER,
	token.GTE:           LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.STAR:          PRODUCT,
	token.FSLASH:        PRODUCT,
	token.LPAREN:        CALL,
	token.DOT:           MEMBERACCESS,
	token.LSQB:          INDEX,
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

// Parses the entire program and builds and returns the AST representation
// of the entire program.
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

func (p *Parser) ParseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.CONST:
		return p.parseConstStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFunctionStatement()
	case token.CLASS:
		return p.parseClassDeclaration()
	case token.FOR:
		return p.parseForInStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// ---------------------------------------------------------------------------
// Statement parsers
// ---------------------------------------------------------------------------

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

	stmt.Type = p.parseTypeHint()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpr(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		// skip tokens until semicolon (for now)
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	stmt.Type = p.parseTypeHint()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpr(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

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

	// skip next tokens until we hit semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	/*
		Parses an expression statement
		Examples
		## Prefix expressions

		-4 -variable_name, !func()

		## Infix expressions

		x + - * / y
		x == y
		x <= y
		x >= y
		x < y
		x > y
		x != y
	*/
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpr(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	/*
		Parses a block statement of the format { ... }

		Current Token Head -> "{"
		reads statements one by one in the group
		returns the list of statements

		Current Token Head After Return -> "}"
	*/
	if !p.curTokenIs(token.LBRACE) {
		p.addParseError(fmt.Sprintf("expected '{' at beginning of block statement, got %s", p.curToken.Literal))
		return nil
	}

	block := &ast.BlockStatement{
		Token: p.curToken, // consume "{"
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

// ---------------------------------------------------------------------------
// Loop statement parsers
// ---------------------------------------------------------------------------

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
	} else {
		p.addParseError("atleast one loop variable needed")
		return nil
	}

	// consume `in`
	p.nextToken()
	stmt.Iterable = p.parseExpr(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

// ---------------------------------------------------------------------------
// Function & class statement parsers
// ---------------------------------------------------------------------------

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

	// -- : returnType
	stmt.ReturnType = p.parseTypeHint()

	// {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
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

		typeHint := p.parseTypeHint()
		if typeHint != nil {
			argumentTypes[ident.Value] = typeHint
		}

		p.nextToken()

		if p.curTokenIs(token.ASSIGN) {
			p.nextToken()
			defaultValue := p.parseExpr(LOWEST)
			argumentDefaults[ident.Value] = defaultValue
			p.nextToken()
		}

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return idents, argumentTypes, argumentDefaults
}

func (p *Parser) parseClassDeclaration() *ast.ClassDeclStatement {
	stmt := &ast.ClassDeclStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
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
			return nil
		}
	}

	// }
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Fields, stmt.Methods = p.parseClassBlockStatements()
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
		} else if p.curTokenIs(token.FUNCTION) {
			meth := p.parseFunctionStatement()
			if meth != nil {
				methods = append(methods, meth)
			}
		} else {
			p.addParseError(fmt.Sprintf("invalid keyword %s inside class definition", p.curToken.Literal))
			return fields, methods
		}

		p.nextToken()
	}

	return fields, methods
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

	stmt.Type = p.parseTypeHint()
	if stmt.Type == nil {
		p.addParseError(fmt.Sprintf("expected `: <type>` after identifier %s", stmt.Name.Value))
		return nil
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

// ---------------------------------------------------------------------------
// Expression parsers
// ---------------------------------------------------------------------------

// Parses an expression based on the current token and the precedence of the expression.
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

// Parses a prefix expression, see `registerPrefixFunc` for examples of prefix expressions.
func (p *Parser) parsePrefixExpr() ast.Expression {
	expr := &ast.PrefixExpr{
		Token: p.curToken,
		Op:    p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpr(PREFIX)

	return expr
}

// Parses an infix expression, see `registerInfixFunc` for examples of infix expressions.
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

func (p *Parser) parseAssignExpr(left ast.Expression) ast.Expression {
	stmt := &ast.AssignExpression{Token: p.curToken}

	if left != nil {
		stmt.Name = left
	} else {
		p.addParseError("invalid left operand: null expression")
		return nil
	}

	stmt.Operator = p.curToken.Literal
	p.nextToken()
	stmt.Value = p.parseExpr(LOWEST)

	return stmt
}

func (p *Parser) parseIfExpr() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.curToken, // consume 'if'
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expr.Condition = p.parseGroupedExpr()
	if expr.Condition == nil {
		return nil
	}

	expr.ThenBranch = p.parseBranchBody()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		expr.ElseBranch = p.parseBranchBody()
	}

	return expr
}

func (p *Parser) parseBranchBody() ast.Expression {
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		return p.parseBlockStatement()
	}
	p.nextToken()
	return p.parseExpr(LOWEST)
}

func (p *Parser) parseCallExpr(function ast.Expression) ast.Expression {
	stmt := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}

	stmt.Arguments = p.parseCallArguments()
	return stmt
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
		return nil
	}

	return args
}

func (p *Parser) parseMemberAccessExpr(left ast.Expression) ast.Expression {
	expr := &ast.MemberAccessExpr{
		Token:  p.curToken,
		Object: left,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	expr.Member = p.parseIdentifier().(*ast.Identifier)
	return expr
}

func (p *Parser) parseIndexExpr(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Name:  left,
	}

	p.nextToken() // consume "["
	expr.Index = p.parseExpr(LOWEST)
	if expr.Index == nil {
		p.addParseError("invalid index expression")
		return nil
	}

	if !p.expectPeek(token.RSQB) {
		return nil
	}

	return expr
}

func (p *Parser) parseGroupedExpr() ast.Expression {
	p.nextToken()
	exp := p.parseExpr(LOWEST)
	if exp == nil {
		return nil
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// ---------------------------------------------------------------------------
// Literal parsers
// ---------------------------------------------------------------------------

// Parses an identifier literal into an ast.Identifier node.
// Identifiers can be variable names, function names, class names, etc.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// Parses an integer literal into an ast.IntegerLiteral node.
func (p *Parser) parseIntegerExpr() ast.Expression {
	il := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addParseError(fmt.Sprintf("could not parse %s as integer: %v", p.curToken.Literal, err))
		return nil
	}
	il.Value = value
	return il
}

// Parses a float literal into an ast.FloatLiteral node.
func (p *Parser) parseFloatExpr() ast.Expression {
	fl := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addParseError(fmt.Sprintf("could not parse %s as float: %v", p.curToken.Literal, err))
		return nil
	}

	fl.Value = value
	return fl
}

// Parses a boolean literal into an ast.Boolean node.
// Boolean literals are `true` or `false`.
func (p *Parser) parseBooleanExpr() ast.Expression {
	b := &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}

	return b
}

// parses a string literal into an ast.StringLiteral node.
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.NullLiteral{
		Token: p.curToken,
	}
}

func (p *Parser) parseTypeExpression() ast.TypeExpression {
	switch p.curToken.Type {
	case token.ARRAY:
		return p.parseArrayType()
	case token.SET:
		return p.parseSetType()
	case token.MAP:
		return p.parseMapType()
	case token.TUPLE:
		return p.parseTupleType()
	case token.UNION:
		return p.parseUnionType()
	case token.IDENT:
		return &ast.SimpleType{
			Token: p.curToken,
			Name:  p.curToken.Literal,
		}
	default:
		p.addParseError(fmt.Sprintf("invalid type specifier: %s", p.curToken.Literal))
		return nil
	}
}

// parses generic type arguments for generic types
// expected is the number of generic type arguments expected, used for error handling
func (p *Parser) _parseGenericTypeParams(expected ...int) ([]ast.TypeExpression, error) {
	if !p.expectPeek(token.LT) {
		return nil, fmt.Errorf("expected '<' to start generic type parameters, got %s", p.peekToken.Literal)
	}

	params := make([]ast.TypeExpression, 0)

	for !p.peekTokenIs(token.GT) {
		p.nextToken()
		typeParam := p.parseTypeExpression()
		if typeParam == nil {
			return nil, fmt.Errorf("invalid type expression in generic type parameters")
		}
		params = append(params, typeParam)
	}

	if len(expected) != 0 && len(params) != expected[0] {
		return nil, fmt.Errorf("expected %d generic type parameters, got %d", expected, len(params))
	}

	if !p.expectPeek(token.GT) {
		return nil, fmt.Errorf("expected '>' to end generic type parameters, got %s", p.peekToken.Literal)
	}

	return params, nil
}

func (p *Parser) parseArrayType() ast.TypeExpression {
	curToken := p.curToken

	if !p.curTokenIs(token.ARRAY) {
		return nil
	}

	elem, err := p._parseGenericTypeParams(1)
	if err != nil {
		p.addParseError(err.Error())
		return nil
	}

	return &ast.ArrayType{
		Token: curToken,
		Elem:  elem[0],
	}
}

func (p *Parser) parseSetType() ast.TypeExpression {
	curToken := p.curToken

	if !p.curTokenIs(token.SET) {
		return nil
	}

	params, err := p._parseGenericTypeParams(1)
	if err != nil {
		p.addParseError(err.Error())
		return nil
	}

	return &ast.SetType{
		Token: curToken,
		Elem:  params[0],
	}
}

func (p *Parser) parseMapType() ast.TypeExpression {
	curToken := p.curToken

	if !p.curTokenIs(token.MAP) {
		return nil
	}

	params, err := p._parseGenericTypeParams(2)
	if err != nil {
		p.addParseError(err.Error())
		return nil
	}

	return &ast.MapType{
		Token: curToken,
		Key:   params[0],
		Value: params[1],
	}
}

func (p *Parser) parseTupleType() ast.TypeExpression {
	curToken := p.curToken

	if !p.curTokenIs(token.TUPLE) {
		return nil
	}

	params, err := p._parseGenericTypeParams()
	if err != nil {
		p.addParseError(err.Error())
		return nil
	}

	return &ast.TupleType{
		Token: curToken,
		Elems: params,
	}
}

func (p *Parser) parseUnionType() ast.TypeExpression {
	curToken := p.curToken

	if !p.curTokenIs(token.UNION) {
		return nil
	}

	params, err := p._parseGenericTypeParams()
	if err != nil {
		p.addParseError(err.Error())
		return nil
	}

	return &ast.UnionType{
		Token: curToken,
		Types: params,
	}
}

// ---------------------------------------------------------------------------
// Utilities & helpers
// ---------------------------------------------------------------------------

// Initializes the parser's prefix and infix parse functions
func initParseFuncs(p *Parser) {
	// Prefix fns
	p.registerPrefixFunc(token.IDENT, p.parseIdentifier)
	p.registerPrefixFunc(token.INT, p.parseIntegerExpr)
	p.registerPrefixFunc(token.FLOAT, p.parseFloatExpr)
	p.registerPrefixFunc(token.STRING, p.parseStringLiteral)
	p.registerPrefixFunc(token.NULL, p.parseNull)
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
	p.registerInfixFunc(token.ASSIGN, p.parseAssignExpr)
	p.registerInfixFunc(token.PLUS_ASSIGN, p.parseAssignExpr)
	p.registerInfixFunc(token.MINUS_ASSIGN, p.parseAssignExpr)
	p.registerInfixFunc(token.STAR_ASSIGN, p.parseAssignExpr)
	p.registerInfixFunc(token.FSLASH_ASSIGN, p.parseAssignExpr)
	p.registerInfixFunc(token.LT, p.parseInfixExpr)
	p.registerInfixFunc(token.GT, p.parseInfixExpr)
	p.registerInfixFunc(token.LTE, p.parseInfixExpr)
	p.registerInfixFunc(token.GTE, p.parseInfixExpr)
	p.registerInfixFunc(token.LPAREN, p.parseCallExpr)
	p.registerInfixFunc(token.DOT, p.parseMemberAccessExpr)
	p.registerInfixFunc(token.LSQB, p.parseIndexExpr)
}

func (p *Parser) registerPrefixFunc(token token.TokenType, fn TPrefixParseFunc) {
	p.prefixParseFuncs[token] = fn
}

func (p *Parser) registerInfixFunc(token token.TokenType, fn TInfixParseFunc) {
	p.infixParseFuncs[token] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
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

func (p *Parser) parseTypeHint() *ast.Identifier {
	if !p.peekTokenIs(token.COLON) {
		return nil
	}
	p.nextToken() // consume ':'

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}
