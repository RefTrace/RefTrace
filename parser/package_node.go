package parser

// PackageNode represents a package in the AST.
type PackageNode struct {
	AnnotatedNode
	name string
}

// NewPackageNode creates a new PackageNode with the given name.
func NewPackageNode(name string) *PackageNode {
	return &PackageNode{
		name: name,
	}
}

// GetName returns the name of the package.
func (p *PackageNode) GetName() string {
	return p.name
}

// GetText returns the text display of this package definition.
func (p *PackageNode) GetText() string {
	return "package " + p.name
}
