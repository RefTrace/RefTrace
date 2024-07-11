package parser

// MixinNode represents a mixin which can be applied to any ClassNode to implement mixins
type MixinNode struct {
	*ClassNode
}

// EmptyMixinNodeArray is an empty array of MixinNode pointers
var EmptyMixinNodeArray = []*MixinNode{}

// NewMixinNode creates a new MixinNode with the given parameters
func NewMixinNode(name string, modifiers int, superType *ClassNode) *MixinNode {
	return NewMixinNodeWithInterfaces(name, modifiers, superType, EmptyClassNodeArray)
}

// NewMixinNodeWithInterfaces creates a new MixinNode with the given parameters including interfaces
func NewMixinNodeWithInterfaces(name string, modifiers int, superType *ClassNode, interfaces []*ClassNode) *MixinNode {
	return &MixinNode{
		ClassNode: NewClassNodeWithInterfaces(name, modifiers, superType, interfaces, EmptyMixinNodeArray),
	}
}
