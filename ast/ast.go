package ast

import (
	"bytes"
	"strings"

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
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var s bytes.Buffer

	s.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		s.WriteString(rs.ReturnValue.String())
	}

	s.WriteString(";")

	return s.String()
}

type ForInStatement struct {
	For      token.Token
	Targets  []*Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fi *ForInStatement) TokenLiteral() string { return fi.For.Literal }
func (fi *ForInStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fi.For.Literal)
	targets := []string{}

	for _, t := range fi.Targets {
		targets = append(targets, t.String())
	}

	out.WriteString(" " + strings.Join(targets, ", "))
	out.WriteString(" " + "in")
	out.WriteString(" " + fi.Iterable.String())
	out.WriteString(" {")
	out.WriteString(fi.Body.String())
	out.WriteString(" }")

	return out.String()
}
func (fi *ForInStatement) statementNode() {}

type ForStatement struct {
	For    token.Token
	Init   Statement
	Cond   Expression
	Update Expression
	Body   *BlockStatement
}

func (fs *ForStatement) TokenLiteral() string { return fs.For.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.For.Literal)
	out.WriteString(" (" + fs.Init.String() + "; ")
	out.WriteString(" " + fs.Cond.String() + "; ")
	out.WriteString(" " + fs.Update.String() + ") ")
	out.WriteString("{")
	out.WriteString(fs.Body.String())
	out.WriteString("}")

	return out.String()
}
func (fs *ForStatement) statementNode() {}

type ClassVarStatement struct {
	Name    *Identifier
	Type    *Identifier
	Default Expression
}

func (cv *ClassVarStatement) statementNode()       {}
func (cv *ClassVarStatement) TokenLiteral() string { return cv.Name.String() }
func (cv *ClassVarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cv.Name.String())
	if cv.Type != nil {
		out.WriteString(": " + cv.Type.String())
	}

	if cv.Default != nil {
		out.WriteString(" = " + cv.Default.String())
	}

	return out.String()
}

type ClassDeclStatement struct {
	Token      token.Token // holds the class name
	Name       *Identifier
	SuperClass *Identifier
	Fields     []*ClassVarStatement
	Methods    []*FunctionStatement
}

func (cs *ClassDeclStatement) TokenLiteral() string { return cs.Token.Literal } // `class`
func (cs *ClassDeclStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral())
	out.WriteString(" " + cs.Name.String())
	if cs.SuperClass != nil {
		out.WriteString("(" + cs.SuperClass.String() + ")")
	}

	out.WriteString(" {")
	for _, field := range cs.Fields {
		out.WriteString(field.String() + "\n")
	}

	for _, meth := range cs.Methods {
		out.WriteString(meth.String() + "\n")
	}
	out.WriteString("}")

	return out.String()
}
func (cs *ClassDeclStatement) statementNode() {}

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

type FunctionStatement struct {
	Token token.Token

	Name       *Identifier
	Params     []*Identifier
	ParamType  map[string]*Identifier
	Defaults   map[string]Expression
	Body       *BlockStatement
	ReturnType *Identifier
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" " + fs.Name.TokenLiteral())
	out.WriteString("(")

	params := []string{}
	for _, p := range fs.Params {
		params = append(params, p.String())
	}

	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")

	if fs.ReturnType != nil {
		out.WriteString(": " + fs.ReturnType.String())
	}

	out.WriteString("{")
	out.WriteString(fs.Body.String())
	out.WriteString("}")

	return out.String()
}

type IfExpression struct {
	/*
		 * Rule
		* IfExpression
		*  :== if "(" <Expression> ")" <Statement> else <Statement>
	*/

	Token      token.Token
	Condition  Expression
	ThenBranch *BlockStatement
	ElseBranch *BlockStatement
}

func (is *IfExpression) TokenLiteral() string { return is.Token.Literal }
func (is *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.ThenBranch.String())

	if is.ElseBranch != nil {
		out.WriteString("else")
		out.WriteString(is.ElseBranch.String())
	}

	return out.String()
}
func (is *IfExpression) expressionNode() {}

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
