package parser

// Variable interface marks an AstNode as a Variable. Typically these are
// VariableExpression, FieldNode, PropertyNode and Parameter
type Variable interface {
	// Name returns the name of the variable.
	Name() string

	// Type returns the type of the variable.
	Type() *ClassNode

	// OriginType returns the type before wrapping primitives type of the variable.
	OriginType() *ClassNode

	// InitialExpression returns the expression used to initialize the variable or nil if there
	// is no initialization.
	GetInitialExpression() Expression

	// HasInitialExpression returns true if there is an initialization expression.
	HasInitialExpression() bool

	// IsClosureSharedVariable returns true if the variable is shared in a closure.
	IsClosureSharedVariable() bool

	// SetClosureSharedVariable sets whether the variable is shared in a closure.
	SetClosureSharedVariable(bool)

	// IsInStaticContext returns true if this variable is used in a static context.
	// A static context is any static initializer block, when this variable
	// is declared as static or when this variable is used in a static method
	IsInStaticContext() bool

	// IsDynamicTyped returns true if the variable is dynamically typed.
	IsDynamicTyped() bool

	// Modifiers returns the modifiers of the variable.
	Modifiers() int

	// IsFinal returns true if the variable is final.
	IsFinal() bool

	// IsPrivate returns true if the variable is private.
	IsPrivate() bool

	// IsProtected returns true if the variable is protected.
	IsProtected() bool

	// IsPublic returns true if the variable is public.
	IsPublic() bool

	// IsStatic returns true if the variable is static.
	IsStatic() bool

	// IsVolatile returns true if the variable is volatile.
	IsVolatile() bool
}

// DefaultVariable provides default implementations for some methods of the Variable interface
type DefaultVariable struct{}

func (DefaultVariable) IsClosureSharedVariable() bool { return false }
func (DefaultVariable) SetClosureSharedVariable(bool) {}

func (v DefaultVariable) Modifiers() int { return ACC_NONE }

func (v DefaultVariable) IsFinal() bool     { return (v.Modifiers() & ACC_FINAL) != 0 }
func (v DefaultVariable) IsPrivate() bool   { return (v.Modifiers() & ACC_PRIVATE) != 0 }
func (v DefaultVariable) IsProtected() bool { return (v.Modifiers() & ACC_PROTECTED) != 0 }
func (v DefaultVariable) IsPublic() bool    { return (v.Modifiers() & ACC_PUBLIC) != 0 }
func (v DefaultVariable) IsStatic() bool    { return (v.Modifiers() & ACC_STATIC) != 0 }
func (v DefaultVariable) IsVolatile() bool  { return (v.Modifiers() & ACC_VOLATILE) != 0 }
