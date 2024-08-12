package parser

// InterfaceHelperClassNode represents an inner class defined as helper for an interface.
type InterfaceHelperClassNode struct {
	*InnerClassNode
	callSites []string
}

// NewInterfaceHelperClassNode creates a new InterfaceHelperClassNode with the given parameters.
func NewInterfaceHelperClassNode(outerClass *ClassNode, name string, modifiers int, superClass *ClassNode, callSites []string) *InterfaceHelperClassNode {
	icn := &InterfaceHelperClassNode{
		InnerClassNode: NewInnerClassNodeWithInterfaces(nil, name, modifiers, superClass, []IClassNode{}, []*MixinNode{}),
		callSites:      make([]string, 0),
	}
	icn.SetCallSites(callSites)
	return icn
}

// SetCallSites sets the call sites for this interface helper class.
func (ihn *InterfaceHelperClassNode) SetCallSites(cs []string) {
	if cs != nil {
		ihn.callSites = cs
	} else {
		ihn.callSites = make([]string, 0)
	}
}

// GetCallSites returns the call sites for this interface helper class.
func (ihn *InterfaceHelperClassNode) GetCallSites() []string {
	return ihn.callSites
}
