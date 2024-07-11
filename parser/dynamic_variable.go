package parser

// DynamicVariable represents an implicitly created variable, such as a variable in a script
// that doesn't have an explicit declaration, or the "it" argument to a closure.
type DynamicVariable struct {
	Variable
	name          string
	closureShare  bool
	staticContext bool
}

// NewDynamicVariable creates a new DynamicVariable with the given name and static context
func NewDynamicVariable(name string, context bool) *DynamicVariable {
	return &DynamicVariable{
		name:          name,
		staticContext: context,
	}
}

// GetType returns the type of the variable (always dynamic)
func (d *DynamicVariable) GetType() *ClassNode {
	return dynamicType()
}

// GetName returns the name of the variable
func (d *DynamicVariable) GetName() string {
	return d.name
}

// GetInitialExpression returns the initial expression (always nil for DynamicVariable)
func (d *DynamicVariable) GetInitialExpression() Expression {
	return nil
}

// HasInitialExpression returns whether the variable has an initial expression (always false)
func (d *DynamicVariable) HasInitialExpression() bool {
	return false
}

// IsInStaticContext returns whether the variable is in a static context
func (d *DynamicVariable) IsInStaticContext() bool {
	return d.staticContext
}

// IsDynamicTyped returns whether the variable is dynamically typed (always true)
func (d *DynamicVariable) IsDynamicTyped() bool {
	return true
}

// IsClosureSharedVariable returns whether the variable is shared in a closure
func (d *DynamicVariable) IsClosureSharedVariable() bool {
	return d.closureShare
}

// SetClosureSharedVariable sets whether the variable is shared in a closure
func (d *DynamicVariable) SetClosureSharedVariable(inClosure bool) {
	d.closureShare = inClosure
}

// GetModifiers returns the modifiers of the variable (always 0 for DynamicVariable)
func (d *DynamicVariable) GetModifiers() int {
	return 0
}

// GetOriginType returns the origin type of the variable (same as GetType for DynamicVariable)
func (d *DynamicVariable) OriginType() *ClassNode {
	return d.Type()
}
