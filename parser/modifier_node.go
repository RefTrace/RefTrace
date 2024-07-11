package parser

import (
	"fmt"
	"hash/fnv"
)

// ModifierNode represents a modifier
type ModifierNode struct {
	ASTNode
	Type           int
	Opcode         int
	Text           string
	AnnotationNode *AnnotationNode
	Repeatable     bool
}

const (
	AnnotationType = -999
)

var ModifierOpcodeMap = map[int]int{
	AnnotationType: 0,
	DEF:            0,
	VAR:            0,

	NATIVE:       ACC_NATIVE,
	SYNCHRONIZED: ACC_SYNCHRONIZED,
	TRANSIENT:    ACC_TRANSIENT,
	VOLATILE:     ACC_VOLATILE,

	PUBLIC:     ACC_PUBLIC,
	PROTECTED:  ACC_PROTECTED,
	PRIVATE:    ACC_PRIVATE,
	STATIC:     ACC_STATIC,
	ABSTRACT:   ACC_ABSTRACT,
	SEALED:     0,
	NON_SEALED: 0,
	FINAL:      ACC_FINAL,
	STRICTFP:   ACC_STRICT,
	DEFAULT:    0,
}

func NewModifierNode(modType int) *ModifierNode {
	opcode, ok := ModifierOpcodeMap[modType]
	if !ok {
		panic(fmt.Sprintf("Unsupported modifier type: %d", modType))
	}

	return &ModifierNode{
		Type:       modType,
		Opcode:     opcode,
		Repeatable: modType == AnnotationType,
	}
}

func NewModifierNodeWithText(modType int, text string) *ModifierNode {
	mn := NewModifierNode(modType)
	mn.Text = text
	return mn
}

func NewModifierNodeWithAnnotation(annotationNode *AnnotationNode, text string) *ModifierNode {
	if annotationNode == nil {
		panic("annotationNode cannot be nil")
	}

	mn := NewModifierNodeWithText(AnnotationType, text)
	mn.AnnotationNode = annotationNode
	return mn
}

func (mn *ModifierNode) IsModifier() bool {
	return !mn.IsAnnotation() && !mn.IsDef()
}

func (mn *ModifierNode) IsVisibilityModifier() bool {
	return mn.Type == PUBLIC || mn.Type == PROTECTED || mn.Type == PRIVATE
}

func (mn *ModifierNode) IsNonVisibilityModifier() bool {
	return mn.IsModifier() && !mn.IsVisibilityModifier()
}

func (mn *ModifierNode) IsAnnotation() bool {
	return mn.Type == AnnotationType
}

func (mn *ModifierNode) IsDef() bool {
	return mn.Type == DEF || mn.Type == VAR
}

func (mn *ModifierNode) GetType() int {
	return mn.Type
}

func (mn *ModifierNode) GetOpcode() int {
	return mn.Opcode
}

func (mn *ModifierNode) IsRepeatable() bool {
	return mn.Repeatable
}

func (mn *ModifierNode) GetText() string {
	return mn.Text
}

func (mn *ModifierNode) GetAnnotationNode() *AnnotationNode {
	return mn.AnnotationNode
}

func (mn *ModifierNode) Equals(o interface{}) bool {
	if mn == o {
		return true
	}
	if o == nil || mn.GetType() != o.(*ModifierNode).GetType() {
		return false
	}
	that := o.(*ModifierNode)
	return mn.Type == that.Type &&
		mn.Text == that.Text &&
		mn.AnnotationNode == that.AnnotationNode
}

func (mn *ModifierNode) HashCode() int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d%s%v", mn.Type, mn.Text, mn.AnnotationNode)))
	return int(h.Sum32())
}

func (mn *ModifierNode) String() string {
	return mn.Text
}
