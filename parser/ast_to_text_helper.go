package parser

import (
	"fmt"
	"strings"
)

// GetClassText returns a string representation of a ClassNode
func GetClassText(node *ClassNode) string {
	if node == nil {
		return "<unknown>"
	}
	return node.ToString(false)
}

// GetParameterText returns a string representation of a Parameter
func GetParameterText(node *Parameter) string {
	if node == nil {
		return "<unknown>"
	}

	name := node.GetName()
	if name == "" {
		name = "<unknown>"
	}
	typeStr := GetClassText(node.GetType())

	text := fmt.Sprintf("%s %s", typeStr, name)
	if node.HasInitialExpression() {
		text += fmt.Sprintf(" = %s", node.GetInitialExpression().Text())
	}
	return text
}

// GetParametersText returns a string representation of a slice of Parameters
func GetParametersText(parameters []*Parameter) string {
	if len(parameters) == 0 {
		return ""
	}
	var result []string
	for _, param := range parameters {
		result = append(result, GetParameterText(param))
	}
	return strings.Join(result, ", ")
}

// GetThrowsClauseText returns a string representation of exception classes
func GetThrowsClauseText(exceptions []*ClassNode) string {
	if len(exceptions) == 0 {
		return ""
	}
	var result []string
	for _, exception := range exceptions {
		result = append(result, GetClassText(exception))
	}
	return " throws " + strings.Join(result, ", ")
}

// GetModifiersText returns a string representation of modifiers
func GetModifiersText(modifiers int) string {
	var result []string

	if modifiers&ACC_PRIVATE != 0 {
		result = append(result, "private")
	}
	if modifiers&ACC_PROTECTED != 0 {
		result = append(result, "protected")
	}
	if modifiers&ACC_PUBLIC != 0 {
		result = append(result, "public")
	}
	if modifiers&ACC_STATIC != 0 {
		result = append(result, "static")
	}
	if modifiers&ACC_ABSTRACT != 0 {
		result = append(result, "abstract")
	}
	if modifiers&ACC_FINAL != 0 {
		result = append(result, "final")
	}
	if modifiers&ACC_INTERFACE != 0 {
		result = append(result, "interface")
	}
	if modifiers&ACC_NATIVE != 0 {
		result = append(result, "native")
	}
	if modifiers&ACC_SYNCHRONIZED != 0 {
		result = append(result, "synchronized")
	}
	if modifiers&ACC_TRANSIENT != 0 {
		result = append(result, "transient")
	}
	if modifiers&ACC_VOLATILE != 0 {
		result = append(result, "volatile")
	}

	return strings.Join(result, " ")
}

// Constants for modifiers (you may need to adjust these values based on your ASM opcodes)
const (
	ACC_PUBLIC       = 0x0001
	ACC_PRIVATE      = 0x0002
	ACC_PROTECTED    = 0x0004
	ACC_STATIC       = 0x0008
	ACC_FINAL        = 0x0010
	ACC_SYNCHRONIZED = 0x0020
	ACC_VOLATILE     = 0x0040
	ACC_TRANSIENT    = 0x0080
	ACC_NATIVE       = 0x0100
	ACC_INTERFACE    = 0x0200
	ACC_ABSTRACT     = 0x0400
)
