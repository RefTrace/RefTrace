package parser

import (
	"strings"
)

// MethodDescriptorWithoutReturnType returns the method's descriptor without the return type
func MethodDescriptorWithoutReturnType(mNode *MethodNode) string {
	var sb strings.Builder
	sb.WriteString(mNode.name)
	sb.WriteString(":")
	for _, p := range mNode.parameters {
		sb.WriteString(FormatTypeName(p.Type()))
		sb.WriteString(",")
	}
	return sb.String()
}

// MethodDescriptor returns the method's full descriptor
func MethodDescriptor(mNode *MethodNode, pretty bool) string {
	name := mNode.name
	if pretty {
		pretty = strings.Contains(name, " ")
	}

	var sb strings.Builder
	sb.WriteString(FormatTypeName(mNode.returnType))
	sb.WriteString(" ")
	if pretty {
		sb.WriteString("\"")
	}
	sb.WriteString(name)
	if pretty {
		sb.WriteString("\"")
	}
	sb.WriteString("(")
	for i, p := range mNode.parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(FormatTypeName(p.Type()))
	}
	sb.WriteString(")")
	return sb.String()
}

// GetPropertyName returns the property name for a potential property method
func GetPropertyName(mNode *MethodNode) string {
	name := mNode.name
	if len(name) > 2 {
		switch name[0] {
		case 'g':
			if len(name) > 3 && name[1] == 'e' && name[2] == 't' && len(mNode.parameters) == 0 && !IsVoidMethod(mNode) {
				return Decapitalize(name[3:])
			}
		case 's':
			if len(name) > 3 && name[1] == 'e' && name[2] == 't' && len(mNode.parameters) == 1 {
				return Decapitalize(name[3:])
			}
		case 'i':
			if name[1] == 's' && len(mNode.parameters) == 0 && IsPrimitiveBoolean(mNode.returnType) {
				return Decapitalize(name[2:])
			}
		}
	}
	return ""
}

// GetCodeAsBlock returns the method's code as a BlockStatement
func GetCodeAsBlock(mNode *MethodNode) *BlockStatement {
	if mNode.code == nil {
		return &BlockStatement{}
	}
	if block, ok := mNode.code.(*BlockStatement); ok {
		return block
	}
	return &BlockStatement{statements: []Statement{mNode.code}}
}

// IsGetterCandidate determines if the method is a getter candidate
func IsGetterCandidate(mNode *MethodNode) bool {
	return len(mNode.parameters) == 0 &&
		IsPublic(mNode.modifiers) &&
		!IsStatic(mNode.modifiers) &&
		!IsAbstract(mNode.modifiers) &&
		!IsVoidMethod(mNode)
}

// WithDefaultArgumentMethods returns a new list including methods for default arguments
func WithDefaultArgumentMethods(methods []*MethodNode) []*MethodNode {
	result := make([]*MethodNode, 0, len(methods))

	for _, method := range methods {
		result = append(result, method)

		if !HasDefaultValue(method) {
			continue
		}

		// ... (implementation for default argument methods)
		// This part would require more context about how default arguments are handled in your Go implementation
	}

	return result
}

// Helper functions (to be implemented based on your specific needs)
func FormatTypeName(cn *ClassNode) string {
	return cn.name // Simplified for this example
}

func Decapitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func IsVoidMethod(mNode *MethodNode) bool {
	return mNode.returnType.name == "void"
}

func IsPrimitiveBoolean(cn *ClassNode) bool {
	return cn.name == "bool"
}

func IsPublic(modifiers int) bool {
	// Implement based on your modifiers system
	return true
}

func IsStatic(modifiers int) bool {
	// Implement based on your modifiers system
	return false
}

func IsAbstract(modifiers int) bool {
	// Implement based on your modifiers system
	return false
}

func HasDefaultValue(mNode *MethodNode) bool {
	// Implement based on your method node structure
	return false
}
