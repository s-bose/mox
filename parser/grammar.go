package parser

import (
	"github.com/s-bose/mox/ast"
	"github.com/s-bose/mox/token"
)

func ParseReturnStatement(p *Parser) *ast.ReturnStatement {
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

func ParseLetStatement(p *Parser) *ast.LetStatement {
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
