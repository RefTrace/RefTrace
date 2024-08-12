package parser

var _ IClassNode = (*InnerClassNode)(nil)

// InnerClassNode represents an inner class definition.
type InnerClassNode struct {
	*ClassNode
	outerClass IClassNode
	scope      *VariableScope
	anonymous  bool
}

// NewInnerClassNode creates a new InnerClassNode with the given parameters.
func NewInnerClassNode(outerClass IClassNode, name string, modifiers int, superClass IClassNode) *InnerClassNode {
	return NewInnerClassNodeWithInterfaces(outerClass, name, modifiers, superClass, []IClassNode{}, []*MixinNode{})
}

// NewInnerClassNodeWithInterfaces creates a new InnerClassNode with the given parameters, including interfaces and mixins.
func NewInnerClassNodeWithInterfaces(outerClass IClassNode, name string, modifiers int, superClass IClassNode, interfaces []IClassNode, mixins []*MixinNode) *InnerClassNode {
	icn := &InnerClassNode{
		ClassNode:  NewClassNodeWithInterfaces(name, modifiers, superClass, interfaces, mixins),
		outerClass: outerClass,
	}
	if outerClass != nil {
		outerClass.AddInnerClass(icn)
	}
	return icn
}

func (icn *InnerClassNode) IsInnerClass() bool {
	return true
}

func (icn *InnerClassNode) SetEnclosingMethod(m MethodOrConstructorNode) {
	icn.ClassNode.SetEnclosingMethod(m)
}

// GetOuterClass returns the outer class of this inner class.
func (icn *InnerClassNode) GetOuterClass() IClassNode {
	return icn.outerClass
}

// GetOuterMostClass returns the outermost class containing this inner class.
func (icn *InnerClassNode) GetOuterMostClass() *ClassNode {
	var outerClass interface{} = icn.GetOuterClass()
	for outerClass != nil {
		innerClass, isInner := outerClass.(*InnerClassNode)
		if !isInner {
			return outerClass.(*ClassNode)
		}
		outerClass = innerClass.GetOuterClass()
	}
	return outerClass.(*ClassNode)
}

// GetOuterField returns the declared field of the outer class with the given name.
func (icn *InnerClassNode) GetOuterField(name string) *FieldNode {
	return icn.outerClass.GetDeclaredField(name)
}

// GetVariableScope returns the variable scope of this inner class.
func (icn *InnerClassNode) GetVariableScope() *VariableScope {
	return icn.scope
}

// SetVariableScope sets the variable scope for this inner class.
func (icn *InnerClassNode) SetVariableScope(scope *VariableScope) {
	icn.scope = scope
}

// IsAnonymous returns whether this inner class is anonymous.
func (icn *InnerClassNode) IsAnonymous() bool {
	return icn.anonymous
}

// SetAnonymous sets whether this inner class is anonymous.
func (icn *InnerClassNode) SetAnonymous(anonymous bool) {
	icn.anonymous = anonymous
}
