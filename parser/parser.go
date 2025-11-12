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
	s := fmt.Sprintf("expected next token to be %s, got %s instead", p.peekToken.Type, t)
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
	fmt.Printf("parsePrefixExpr() called for token %s\n", p.curToken.Literal)

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
	fmt.Printf("parseExpr() called with precedence %d and curToken %s\n", precedence, p.curToken.Literal)
	prefixFunc := p.prefixParseFuncs[p.curToken.Type]
	if prefixFunc == nil {
		p.addParseError(fmt.Sprintf("no prefix function exists for %s", p.curToken.Type))
		return nil
	}
	fmt.Printf("prefixFunc is %T\n", prefixFunc)

	leftExp := prefixFunc()
	fmt.Printf("leftExp: %s\n", leftExp.String())
	fmt.Printf("precedence: %d, peekPrecedence: %d\n", precedence, p.peekPrecedence())
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

// Parse statements helpers

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
	// skip next tokens until we hit semicolon

	if !p.curTokenIs(token.SEMICOLON) {
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

	for !p.curTokenIs(token.SEMICOLON) {
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

func (p *Parser) ParseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}
