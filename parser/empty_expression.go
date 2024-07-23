package parser

var _ Expression = (*EmptyExpression)(nil)

// EmptyExpression represents a placeholder for an empty expression.
// Empty expressions are used in closures lists like (;).
// During class generation, empty expressions should be ignored
// or replaced with a null value.
type EmptyExpression struct {
	*BaseExpression
}

// NewEmptyExpression creates a new EmptyExpression instance.
func NewEmptyExpression() *EmptyExpression {
	return &EmptyExpression{
		BaseExpression: NewBaseExpression(),
	}
}

// Visit implements the Expression interface.
func (e *EmptyExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitEmptyExpression(e)
}

// INSTANCE is an immutable singleton that is recommended for use when source range
// or any other occurrence-specific metadata is not needed.
var INSTANCE = NewEmptyExpression()

// Ensure INSTANCE methods do nothing when attempting to modify

func (e *EmptyExpression) SetColumnNumber(n int) {
	if e != INSTANCE {
		e.BaseExpression.SetColumnNumber(n)
	}
}

func (e *EmptyExpression) SetLastColumnNumber(n int) {
	if e != INSTANCE {
		e.BaseExpression.SetLastColumnNumber(n)
	}
}

func (e *EmptyExpression) SetLastLineNumber(n int) {
	if e != INSTANCE {
		e.BaseExpression.SetLastLineNumber(n)
	}
}

func (e *EmptyExpression) SetLineNumber(n int) {
	if e != INSTANCE {
		e.BaseExpression.SetLineNumber(n)
	}
}

func (e *EmptyExpression) SetMetaDataMap(meta map[interface{}]interface{}) {
	if e != INSTANCE {
		e.BaseExpression.SetMetaDataMap(meta)
	}
}

func (e *EmptyExpression) SetSourcePosition(node ASTNode) {
	if e != INSTANCE {
		e.BaseExpression.SetSourcePosition(node)
	}
}

func (e *EmptyExpression) AddAnnotation(node *ClassNode) {
	if e != INSTANCE {
		e.BaseExpression.AddAnnotation(node)
	}
}

func (e *EmptyExpression) SetDeclaringClass(node *ClassNode) {
	if e != INSTANCE {
		e.BaseExpression.SetDeclaringClass(node)
	}
}

func (e *EmptyExpression) SetHasNoRealSourcePosition(b bool) {
	if e != INSTANCE {
		e.BaseExpression.SetHasNoRealSourcePosition(b)
	}
}

func (e *EmptyExpression) SetSynthetic(b bool) {
	if e != INSTANCE {
		e.BaseExpression.SetSynthetic(b)
	}
}

func (e *EmptyExpression) SetType(node *ClassNode) {
	if e != INSTANCE {
		e.BaseExpression.SetType(node)
	}
}

// Implement Expression interface methods
func (e *EmptyExpression) GetType() *ClassNode {
	return e.BaseExpression.GetType()
}

// ASTNode interface methods are already implemented by BaseExpression
