package parser

import (
	"fmt"
	"strings"
)

// AnnotationNode represents an annotation which can be attached to interfaces, classes, methods, fields, parameters, and other places.
type AnnotationNode struct {
	BaseASTNode
	classNode        *ClassNode
	members          map[string]Expression
	runtimeRetention bool
	sourceRetention  bool
	classRetention   bool
	allowedTargets   int
}

const (
	ConstructorTarget     = 1 << 1
	MethodTarget          = 1 << 2
	FieldTarget           = 1 << 3
	ParameterTarget       = 1 << 4
	LocalVariableTarget   = 1 << 5
	AnnotationTarget      = 1 << 6
	PackageTarget         = 1 << 7
	TypeParameterTarget   = 1 << 8
	TypeUseTarget         = 1 << 9
	RecordComponentTarget = 1 << 10
	TypeTarget            = 1 + AnnotationTarget
)

const AllTargets = TypeTarget | ConstructorTarget | MethodTarget |
	FieldTarget | ParameterTarget | LocalVariableTarget | AnnotationTarget |
	PackageTarget | TypeParameterTarget | TypeUseTarget | RecordComponentTarget

// NewAnnotationNode creates a new AnnotationNode
func NewAnnotationNode(classNode *ClassNode) *AnnotationNode {
	return &AnnotationNode{
		classNode:      classNode,
		allowedTargets: AllTargets,
	}
}

// GetClassNode returns the ClassNode of this annotation
func (an *AnnotationNode) GetClassNode() *ClassNode {
	return an.classNode
}

// GetMembers returns all members of this annotation
func (an *AnnotationNode) GetMembers() map[string]Expression {
	if an.members == nil {
		return make(map[string]Expression)
	}
	return an.members
}

// GetMember returns a specific member of this annotation
func (an *AnnotationNode) GetMember(name string) Expression {
	if an.members == nil {
		return nil
	}
	return an.members[name]
}

// AddMember adds a new member to this annotation
func (an *AnnotationNode) AddMember(name string, value Expression) error {
	an.assertMembers()
	if _, exists := an.members[name]; exists {
		return fmt.Errorf("annotation member %s has already been added", name)
	}
	an.members[name] = value
	return nil
}

// SetMember sets a member of this annotation
func (an *AnnotationNode) SetMember(name string, value Expression) {
	an.assertMembers()
	an.members[name] = value
}

func (an *AnnotationNode) assertMembers() {
	if an.members == nil {
		an.members = make(map[string]Expression)
	}
}

// IsBuiltIn returns whether this annotation is built-in
func (an *AnnotationNode) IsBuiltIn() bool {
	return false
}

// HasRuntimeRetention returns true if the annotation should be visible at runtime
func (an *AnnotationNode) HasRuntimeRetention() bool {
	return an.runtimeRetention
}

// SetRuntimeRetention sets whether the annotation has runtime retention
func (an *AnnotationNode) SetRuntimeRetention(flag bool) {
	an.runtimeRetention = flag
}

// HasSourceRetention returns true if the annotation is only allowed in sources
func (an *AnnotationNode) HasSourceRetention() bool {
	return an.sourceRetention
}

// SetSourceRetention sets whether the annotation has source retention
func (an *AnnotationNode) SetSourceRetention(flag bool) {
	an.sourceRetention = flag
}

// HasClassRetention returns true if the annotation is written in the bytecode, but not visible at runtime
func (an *AnnotationNode) HasClassRetention() bool {
	if !an.runtimeRetention && !an.sourceRetention {
		return true
	}
	return an.classRetention
}

// SetClassRetention sets whether the annotation has explicit class retention
func (an *AnnotationNode) SetClassRetention(flag bool) {
	an.classRetention = flag
}

// SetAllowedTargets sets the allowed targets for this annotation
func (an *AnnotationNode) SetAllowedTargets(bitmap int) {
	an.allowedTargets = bitmap
}

// IsTargetAllowed checks if a specific target is allowed for this annotation
func (an *AnnotationNode) IsTargetAllowed(target int) bool {
	return (an.allowedTargets & target) == target
}

// TargetToName converts a target constant to its string representation
func TargetToName(target int) string {
	switch target {
	case TypeTarget:
		return "TYPE"
	case ConstructorTarget:
		return "CONSTRUCTOR"
	case MethodTarget:
		return "METHOD"
	case FieldTarget:
		return "FIELD"
	case ParameterTarget:
		return "PARAMETER"
	case LocalVariableTarget:
		return "LOCAL_VARIABLE"
	case AnnotationTarget:
		return "ANNOTATION"
	case PackageTarget:
		return "PACKAGE"
	case TypeParameterTarget:
		return "TYPE_PARAMETER"
	case TypeUseTarget:
		return "TYPE_USE"
	case RecordComponentTarget:
		return "RECORD_COMPONENT"
	default:
		return "unknown target"
	}
}

// String returns a string representation of the AnnotationNode
func (an *AnnotationNode) String() string {
	return fmt.Sprintf("annotationnode[%s]", an.GetText())
}

// GetText returns the text representation of the AnnotationNode
func (an *AnnotationNode) GetText() string {
	var memberText strings.Builder
	if an.members != nil {
		first := true
		for key, value := range an.members {
			if !first {
				memberText.WriteString(", ")
			}
			first = false
			text := value.GetText()
			if listExpr, ok := value.(*ListExpression); ok {
				var result []string
				for _, exp := range listExpr.GetExpressions() {
					if annotationConstExpr, ok := exp.(*AnnotationConstantExpression); ok {
						result = append(result, annotationConstExpr.GetValue().(ASTNode).GetText())
					} else {
						result = append(result, exp.GetText())
					}
				}
				text = fmt.Sprintf("%v", result)
			}
			memberText.WriteString(fmt.Sprintf("%s: %s", key, text))
		}
	}
	return fmt.Sprintf("@%s(%s)", an.classNode.GetText(), memberText.String())
}
