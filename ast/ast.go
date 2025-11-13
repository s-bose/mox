package ast

import (
	"bytes"

	"github.com/s-bose/mox/token"
)

// Interfaces
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// AST
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var s bytes.Buffer

	s.WriteString(ls.TokenLiteral() + " ")
	s.WriteString(ls.Name.String())
	if ls.Type != nil {
		s.WriteString(": " + ls.Type.String())
	}
	s.WriteString(" = ")
	if ls.Value != nil {
		s.WriteString(ls.Value.String())
	}

	s.WriteString(";")
	return s.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var s bytes.Buffer

	s.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		s.WriteString(rs.Value.String())
	}

	s.WriteString(";")

	return s.String()
}

// Expression Nodes
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }
func (il *IntegerLiteral) expressionNode()      {}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (il *FloatLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *FloatLiteral) String() string       { return il.TokenLiteral() }
func (il *FloatLiteral) expressionNode()      {}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.TokenLiteral() }
func (b *Boolean) expressionNode()      {}

type PrefixExpr struct {
	Token token.Token
	Op    string
	Right Expression
}

func (pe *PrefixExpr) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpr) String() string {
	var s bytes.Buffer
	s.WriteString("(")
	s.WriteString(pe.Op)
	s.WriteString(pe.Right.String())
	s.WriteString(")")

	return s.String()
}
func (pe *PrefixExpr) expressionNode() {}

type InfixExpr struct {
	Token token.Token // holds the infix op Token
	Left  Expression
	Op    string // holds the raw string for the operator
	Right Expression
}

func (ie *InfixExpr) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpr) String() string {
	var s bytes.Buffer

	s.WriteString("(")
	s.WriteString(ie.Left.String() + " ")
	s.WriteString(ie.Op)
	s.WriteString(" " + ie.Right.String())
	s.WriteString(")")

	return s.String()
}
func (ie *InfixExpr) expressionNode() {}

type IfStatement struct {
	/*
		 * Rule
		* IfStatement :== if "(" <Expression> ")" <Statement> else <Statement>
	*/

	Token      token.Token
	Condition  Expression
	ThenBranch Statement
	ElseBranch Statement
}

func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.ThenBranch.String())

	if is.ElseBranch != nil {
		out.WriteString(is.ElseBranch.String())
	}

	return out.String()
}
func (is *IfStatement) statementNode()

type ElseStatement struct {
	Token      token.Token
	ThenBranch Statement
}

func (es *ElseStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ElseStatement) String() string {
	var out bytes.Buffer

	out.WriteString("else")
	out.WriteString(es.ThenBranch.String())

	return out.String()
}
func (es *ElseStatement) statementNode() {}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Main Program
type Program struct {
	Statements []Statement
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var s bytes.Buffer

	for _, stmt := range p.Statements {
		s.WriteString(stmt.String())
	}

	return s.String()
}
