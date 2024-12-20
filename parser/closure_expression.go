package parser

// ClosureExpression represents a closure expression such as { statement }
// or { i -> statement } or { i, x, String y ->  statement }
type ClosureExpression struct {
	*BaseExpression
	Parameters    []*Parameter
	Code          Statement
	VariableScope *VariableScope
}

// NewClosureExpression creates a new ClosureExpression
func NewClosureExpression(parameters []*Parameter, code Statement) *ClosureExpression {
	ce := &ClosureExpression{
		BaseExpression: NewBaseExpression(),
		Parameters:     parameters,
		Code:           code,
	}
	ce.SetType(CLOSURE_TYPE.GetPlainNodeReference())
	return ce
}

// GetCode returns the code statement of the closure
func (ce *ClosureExpression) GetCode() Statement {
	return ce.Code
}

// SetCode sets the code statement of the closure
func (ce *ClosureExpression) SetCode(code Statement) {
	ce.Code = code
}

// GetParameters returns an array of zero (for implicit it) or more (when explicit args given) parameters
// or nil otherwise (representing explicit no args)
func (ce *ClosureExpression) GetParameters() []*Parameter {
	return ce.Parameters
}

// IsParameterSpecified returns true if one or more explicit parameters are supplied
func (ce *ClosureExpression) IsParameterSpecified() bool {
	return ce.Parameters != nil && len(ce.Parameters) > 0
}

// GetVariableScope returns the variable scope of the closure
func (ce *ClosureExpression) GetVariableScope() *VariableScope {
	return ce.VariableScope
}

// SetVariableScope sets the variable scope of the closure
func (ce *ClosureExpression) SetVariableScope(variableScope *VariableScope) {
	ce.VariableScope = variableScope
}

// GetText returns a string representation of the closure
func (ce *ClosureExpression) GetText() string {
	return ce.toString("...")
}

// String returns a string representation of the closure
func (ce *ClosureExpression) String() string {
	codeStr := "<null>"
	if ce.Code != nil {
		codeStr = ce.Code.GetText()
	}
	return ce.BaseExpression.GetText() + ce.toString(codeStr)
}

func (ce *ClosureExpression) toString(bodyText string) string {
	if HasImplicitParameter(ce) {
		return "{ " + bodyText + " }"
	}
	return "{ " + GetParametersText(ce.Parameters) + " -> " + bodyText + " }"
}

// TransformExpression transforms the expression
func (ce *ClosureExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	return ce
}

// Visit calls the VisitClosureExpression method of the GroovyCodeVisitor
func (ce *ClosureExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitClosureExpression(ce)
}

// HasImplicitParameter checks if the closure has an implicit parameter
func HasImplicitParameter(ce *ClosureExpression) bool {
	// Implement this function based on the ClosureUtils.hasImplicitParameter logic
	return false // Placeholder implementation
}
