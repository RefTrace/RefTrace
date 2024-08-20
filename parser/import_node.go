package parser

import (
	"fmt"
	"hash/fnv"
)

// ImportNode represents an import statement.
type ImportNode struct {
	*AnnotatedNode
	Type        IClassNode
	Alias       string
	FieldName   string
	PackageName string
	IsStar      bool
	IsStatic    bool
	hashCode    uint32
}

// NewImportNodeType creates an import of a single type
func NewImportNodeType(typ IClassNode, alias string) *ImportNode {
	return &ImportNode{
		AnnotatedNode: NewAnnotatedNode(),
		Type:          typ,
		Alias:         alias,
		IsStar:        false,
		IsStatic:      false,
	}
}

// NewImportNodePackage creates an import of all types in a package
func NewImportNodePackage(packageName string) *ImportNode {
	return &ImportNode{
		AnnotatedNode: NewAnnotatedNode(),
		PackageName:   packageName,
		IsStar:        true,
		IsStatic:      false,
	}
}

// NewImportNodeStatic creates an import of all static members of a type
func NewImportNodeStatic(typ IClassNode) *ImportNode {
	return &ImportNode{
		AnnotatedNode: NewAnnotatedNode(),
		Type:          typ,
		IsStar:        true,
		IsStatic:      true,
	}
}

// NewImportNodeStaticField creates an import of a static field or method of a type
func NewImportNodeStaticField(typ IClassNode, fieldName, alias string) *ImportNode {
	return &ImportNode{
		AnnotatedNode: NewAnnotatedNode(),
		Type:          typ,
		Alias:         alias,
		FieldName:     fieldName,
		IsStatic:      true,
	}
}

// GetText returns the text display of this import
func (in *ImportNode) GetText() string {
	simpleName := in.Alias
	memberName := in.FieldName

	if !in.IsStatic {
		if in.IsStar {
			return fmt.Sprintf("import %s*", in.PackageName)
		} else if simpleName == "" || simpleName == in.Type.GetNameWithoutPackage() {
			return fmt.Sprintf("import %s", in.GetClassName())
		} else {
			return fmt.Sprintf("import %s as %s", in.GetClassName(), simpleName)
		}
	} else {
		if in.IsStar {
			return fmt.Sprintf("import static %s.*", in.GetClassName())
		} else if simpleName == "" || simpleName == memberName {
			return fmt.Sprintf("import static %s.%s", in.GetClassName(), memberName)
		} else {
			return fmt.Sprintf("import static %s.%s as %s", in.GetClassName(), memberName, simpleName)
		}
	}
}

// GetClassName returns the class name
func (in *ImportNode) GetClassName() string {
	if in.Type == nil {
		return ""
	}
	return in.Type.GetName()
}

// SetType sets the type
func (in *ImportNode) SetType(typ IClassNode) {
	in.Type = typ
	in.hashCode = 0
}

// Equals checks if two ImportNodes are equal
func (in *ImportNode) Equals(other interface{}) bool {
	if other, ok := other.(*ImportNode); ok {
		return in.Type == other.Type &&
			in.Alias == other.Alias &&
			in.FieldName == other.FieldName &&
			in.PackageName == other.PackageName &&
			in.IsStar == other.IsStar &&
			in.IsStatic == other.IsStatic
	}
	return false
}

// HashCode returns the hash code for the ImportNode
func (in *ImportNode) HashCode() uint32 {
	if in.hashCode == 0 {
		h := fnv.New32a()
		h.Write([]byte(fmt.Sprintf("%v%s%s%s%v%v", in.Type, in.Alias, in.FieldName, in.PackageName, in.IsStar, in.IsStatic)))
		in.hashCode = h.Sum32()
	}
	return in.hashCode
}

// Visit implements the GroovyCodeVisitor interface
func (in *ImportNode) Visit(visitor GroovyCodeVisitor) {
	// Empty implementation
}
