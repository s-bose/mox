package ast

import (
	"bytes"
	"strings"

	"github.com/s-bose/mox/token"
)

// ---------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------

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

type TypeExpression interface {
	Node
	typeExpressionNode()
}

// ---------------------------------------------------------------------------
// Program – the root AST node
// ---------------------------------------------------------------------------

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var s bytes.Buffer
	for _, stmt := range p.Statements {
		s.WriteString(stmt.String())
	}
	return s.String()
}

// ---------------------------------------------------------------------------
// Literal / Atomic Expressions
// ---------------------------------------------------------------------------

// Identifier

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }

// FloatLiteral

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.TokenLiteral() }

// Boolean

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.TokenLiteral() }

// NullLiteral

type NullLiteral struct {
	Token token.Token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NullLiteral) String() string       { return n.TokenLiteral() }

// StringLiteral

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// --------–-----------------------------------------------------------------
// Composite Types
// --------------------------------------------------------------------------

// SimpleType  e.g. int, string, MyClass

type SimpleType struct {
	Token token.Token
	Name  string
}

func (st *SimpleType) typeExpressionNode()  {}
func (st *SimpleType) TokenLiteral() string { return st.Token.Literal }
func (st *SimpleType) String() string       { return st.Name }

// ArrayType  e.g. array<string>, array<int>, array<MyClass>

type ArrayType struct {
	Token token.Token
	Elem  TypeExpression
}

func (at *ArrayType) typeExpressionNode()  {}
func (at *ArrayType) TokenLiteral() string { return at.Token.Literal }
func (at *ArrayType) String() string {
	var s bytes.Buffer
	s.WriteString("array<")
	s.WriteString(at.Elem.String())
	s.WriteString(">")
	return s.String()
}

// MapType  e.g. map<string, int>, map<string, MyClass>

type MapType struct {
	Token token.Token
	Key   TypeExpression
	Value TypeExpression
}

func (mt *MapType) typeExpressionNode()  {}
func (mt *MapType) TokenLiteral() string { return mt.Token.Literal }
func (mt *MapType) String() string {
	var s bytes.Buffer
	s.WriteString("map<")
	s.WriteString(mt.Key.String())
	s.WriteString(", ")
	s.WriteString(mt.Value.String())
	s.WriteString(">")
	return s.String()
}

// SetType e.g. set<int>, set<string>

type SetType struct {
	Token token.Token
	Elem  TypeExpression
}

func (st *SetType) typeExpressionNode()  {}
func (st *SetType) TokenLiteral() string { return st.Token.Literal }
func (st *SetType) String() string {
	var s bytes.Buffer
	s.WriteString("set<")
	s.WriteString(st.Elem.String())
	s.WriteString(">")
	return s.String()
}

// TupleType e.g. tuple<int, string, MyClass>

type TupleType struct {
	Token token.Token
	Elems []TypeExpression
}

func (tt *TupleType) typeExpressionNode()  {}
func (tt *TupleType) TokenLiteral() string { return tt.Token.Literal }
func (tt *TupleType) String() string {
	var s bytes.Buffer
	elemStrs := []string{}

	for _, e := range tt.Elems {
		elemStrs = append(elemStrs, e.String())
	}

	s.WriteString("tuple<")
	s.WriteString(strings.Join(elemStrs, ", "))
	s.WriteString(">")
	return s.String()
}

// UnionType e.g. union<string, int, MyClass>

type UnionType struct {
	Token token.Token
	Types []TypeExpression
}

func (ut *UnionType) typeExpressionNode()  {}
func (ut *UnionType) TokenLiteral() string { return ut.Token.Literal }
func (ut *UnionType) String() string {
	var s bytes.Buffer
	typeStrs := []string{}
	for _, t := range ut.Types {
		typeStrs = append(typeStrs, t.String())
	}
	s.WriteString("union<")
	s.WriteString(strings.Join(typeStrs, ", "))
	s.WriteString(">")
	return s.String()
}

// ---------------------------------------------------------------------------
// Operator Expressions
// ---------------------------------------------------------------------------

// PrefixExpr  e.g. !x, -5

type PrefixExpr struct {
	Token token.Token
	Op    string
	Right Expression
}

func (pe *PrefixExpr) expressionNode()      {}
func (pe *PrefixExpr) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpr) String() string {
	var s bytes.Buffer
	s.WriteString("(")
	s.WriteString(pe.Op)
	s.WriteString(pe.Right.String())
	s.WriteString(")")
	return s.String()
}

// InfixExpr  e.g. a + b

type InfixExpr struct {
	Token token.Token // holds the infix op Token
	Left  Expression
	Op    string // holds the raw string for the operator
	Right Expression
}

func (ie *InfixExpr) expressionNode()      {}
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

// AssignExpression  e.g. x = 5, x += 1

type AssignExpression struct {
	Token    token.Token // the '=' token
	Name     Expression
	Operator string
	Value    Expression
}

func (as *AssignExpression) expressionNode()      {}
func (as *AssignExpression) TokenLiteral() string { return as.Token.Literal }
func (as *AssignExpression) String() string {
	var s bytes.Buffer
	s.WriteString(as.Name.String())
	s.WriteString(" " + as.Operator + " ")
	s.WriteString(as.Value.String())
	return s.String()
}

// IndexExpression  e.g. arr[0]

type IndexExpression struct {
	Token token.Token
	Name  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var s bytes.Buffer
	s.WriteString("(")
	s.WriteString(ie.Name.String())
	s.WriteString("[")
	s.WriteString(ie.Index.String())
	s.WriteString("])")
	return s.String()
}

// MemberAccessExpr  e.g. obj.field

type MemberAccessExpr struct {
	Token  token.Token // the DOT token
	Object Expression  // expression before the dot
	Member *Identifier // identifier after the dot
}

func (ma *MemberAccessExpr) expressionNode()      {}
func (ma *MemberAccessExpr) TokenLiteral() string { return ma.Token.Literal }
func (ma *MemberAccessExpr) String() string {
	var s bytes.Buffer
	s.WriteString("(")
	s.WriteString(ma.Object.String())
	s.WriteString(".")
	s.WriteString(ma.Member.String())
	s.WriteString(")")
	return s.String()
}

// ---------------------------------------------------------------------------
// Compound Expressions
// ---------------------------------------------------------------------------

// IfExpression
//
//	:== if "(" <Expression> ")" <Statement> else <Statement>

type IfExpression struct {
	Token      token.Token
	Condition  Expression
	ThenBranch Expression
	ElseBranch Expression
}

func (is *IfExpression) expressionNode()      {}
func (is *IfExpression) TokenLiteral() string { return is.Token.Literal }
func (is *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString("( ")
	out.WriteString(is.Condition.String())
	out.WriteString(" )")
	out.WriteString(" ")
	out.WriteString(is.ThenBranch.String())
	if is.ElseBranch != nil {
		out.WriteString("else")
		out.WriteString(is.ElseBranch.String())
	}
	return out.String()
}

// CallExpression  e.g. foo(a, b)

type CallExpression struct {
	Token     token.Token // stores '('
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String() + "(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// BlockStatement  –  { ... }
// Implements both Statement and Expression.

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ---------------------------------------------------------------------------
// Simple Statements
// ---------------------------------------------------------------------------

// ExpressionStatement  –  an expression used as a statement

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// LetStatement  –  let x: int = 5;

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

// ConstStatement  –  const x: int = 5;

type ConstStatement struct {
	Token token.Token
	Name  *Identifier
	Type  *Identifier
	Value Expression
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConstStatement) String() string {
	var s bytes.Buffer
	s.WriteString(cs.TokenLiteral() + " ")
	s.WriteString(cs.Name.String())
	if cs.Type != nil {
		s.WriteString(": " + cs.Type.String())
	}
	s.WriteString(" = ")
	if cs.Value != nil {
		s.WriteString(cs.Value.String())
	}
	s.WriteString(";")
	return s.String()
}

// ReturnStatement  –  return <expr>;

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

// ---------------------------------------------------------------------------
// Loop Statements
// ---------------------------------------------------------------------------

// ForInStatement  –  for x in iterable { ... }

type ForInStatement struct {
	For      token.Token
	Targets  []*Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fi *ForInStatement) statementNode()       {}
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

// ForStatement  –  for (init; cond; update) { ... }

type ForStatement struct {
	For    token.Token
	Init   Statement
	Cond   Expression
	Update Expression
	Body   *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
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

// ---------------------------------------------------------------------------
// Function & Class Statements
// ---------------------------------------------------------------------------

// FunctionStatement  –  fn name(params): returnType { body }

type FunctionStatement struct {
	Token      token.Token
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

// ClassVarStatement  –  a field declaration inside a class body

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

// ClassDeclStatement  –  class Name(Super) { fields; methods }

type ClassDeclStatement struct {
	Token      token.Token // holds the class name
	Name       *Identifier
	SuperClass *Identifier
	Fields     []*ClassVarStatement
	Methods    []*FunctionStatement
}

func (cs *ClassDeclStatement) statementNode()       {}
func (cs *ClassDeclStatement) TokenLiteral() string { return cs.Token.Literal }
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
