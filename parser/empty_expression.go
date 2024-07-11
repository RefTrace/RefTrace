package parser

import (
	"errors"
)

// EmptyExpression represents a placeholder for an empty expression.
// Empty expressions are used in closures lists like (;).
// During class generation, empty expressions should be ignored
// or replaced with a null value.
type EmptyExpression struct {
	Expression
}

// NewEmptyExpression creates a new EmptyExpression instance.
func NewEmptyExpression() *EmptyExpression {
	return &EmptyExpression{}
}

// Visit implements the Expression interface.
func (e *EmptyExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitEmptyExpression(e)
}

// INSTANCE is an immutable singleton that is recommended for use when source range
// or any other occurrence-specific metadata is not needed.
var INSTANCE = &EmptyExpression{}

// Ensure INSTANCE methods throw errors when attempting to modify

func (e *EmptyExpression) SetColumnNumber(n int) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetColumnNumber(n)
	return nil
}

func (e *EmptyExpression) SetLastColumnNumber(n int) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetLastColumnNumber(n)
	return nil
}

func (e *EmptyExpression) SetLastLineNumber(n int) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetLastLineNumber(n)
	return nil
}

func (e *EmptyExpression) SetLineNumber(n int) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetLineNumber(n)
	return nil
}

func (e *EmptyExpression) SetMetaDataMap(meta map[interface{}]interface{}) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetMetaDataMap(meta)
	return nil
}

func (e *EmptyExpression) SetSourcePosition(node ASTNode) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetSourcePosition(node)
	return nil
}

func (e *EmptyExpression) AddAnnotation(node AnnotationNode) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.AddAnnotation(node)
	return nil
}

func (e *EmptyExpression) SetDeclaringClass(node ClassNode) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetDeclaringClass(node)
	return nil
}

func (e *EmptyExpression) SetHasNoRealSourcePosition(b bool) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetHasNoRealSourcePosition(b)
	return nil
}

func (e *EmptyExpression) SetSynthetic(b bool) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetSynthetic(b)
	return nil
}

func (e *EmptyExpression) SetType(node ClassNode) error {
	if e == INSTANCE {
		return errors.New("EmptyExpression.INSTANCE is immutable")
	}
	e.SetType(node)
	return nil
}
