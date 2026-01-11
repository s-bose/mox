package eval

import (
	"github.com/s-bose/mox/ast"
)

type ObjectType string

type IObject interface {
	Type() ObjectType
}

// Object representation of all internal data types
// Contains
// Type()
// Inspect()

// This is a retarded way of doing this but eh

type Evaluator struct {
	// fields for evaluator
	env *Environment
}

func (e *Evaluator) Eval(node ast.Node) IObject {
	// evaluation logic goes here
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &Float{Value: node.Value}

	}

	return nil
}

type Environment struct {
	// fields for environment
}
