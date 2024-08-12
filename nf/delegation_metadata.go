package nf

import "reft-go/parser"

// DelegationMetadata stores the delegation strategy and delegate type of closures.
// As closures can be nested, a delegation metadata may have a parent.
type DelegationMetadata struct {
	parent   *DelegationMetadata
	typ      parser.IClassNode
	strategy int
}

// NewDelegationMetadata creates a new DelegationMetadata instance with a parent
func NewDelegationMetadata(typ parser.IClassNode, strategy int, parent *DelegationMetadata) *DelegationMetadata {
	return &DelegationMetadata{
		strategy: strategy,
		typ:      wrapTypeIfNecessary(typ), // non-primitive
		parent:   parent,
	}
}

// NewDelegationMetadataWithoutParent creates a new DelegationMetadata instance without a parent
func NewDelegationMetadataWithoutParent(typ parser.IClassNode, strategy int) *DelegationMetadata {
	return NewDelegationMetadata(typ, strategy, nil)
}

// GetStrategy returns the delegation strategy
func (d *DelegationMetadata) GetStrategy() int {
	return d.strategy
}

// GetType returns the delegate type
func (d *DelegationMetadata) GetType() parser.IClassNode {
	return d.typ
}

// GetParent returns the parent DelegationMetadata, if any
func (d *DelegationMetadata) GetParent() *DelegationMetadata {
	return d.parent
}
