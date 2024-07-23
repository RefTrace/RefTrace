package parser

// EnumConstantClassNode represents the anonymous inner class for an enum constant.
type EnumConstantClassNode struct {
	*InnerClassNode
}

// NewEnumConstantClassNode creates a new EnumConstantClassNode with the given parameters.
func NewEnumConstantClassNode(outerClass *ClassNode, name string, superClass *ClassNode) *EnumConstantClassNode {
	return &EnumConstantClassNode{
		InnerClassNode: NewInnerClassNode(outerClass, name, ACC_ENUM|ACC_FINAL, superClass),
	}
}
