package parser

import (
	"fmt"
	"hash/fnv"
)

// ModifierNode represents a modifier
type ModifierNode struct {
	*BaseASTNode
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
	AnnotationType:  0,
	GroovyParserDEF: 0,
	GroovyParserVAR: 0,

	GroovyParserNATIVE:       ACC_NATIVE,
	GroovyParserSYNCHRONIZED: ACC_SYNCHRONIZED,
	GroovyParserTRANSIENT:    ACC_TRANSIENT,
	GroovyParserVOLATILE:     ACC_VOLATILE,

	GroovyParserPUBLIC:     ACC_PUBLIC,
	GroovyParserPROTECTED:  ACC_PROTECTED,
	GroovyParserPRIVATE:    ACC_PRIVATE,
	GroovyParserSTATIC:     ACC_STATIC,
	GroovyParserABSTRACT:   ACC_ABSTRACT,
	GroovyParserSEALED:     0,
	GroovyParserNON_SEALED: 0,
	GroovyParserFINAL:      ACC_FINAL,
	GroovyParserSTRICTFP:   ACC_STRICT,
	GroovyParserDEFAULT:    0,
}

func NewModifierNode(modType int) *ModifierNode {
	opcode, ok := ModifierOpcodeMap[modType]
	if !ok {
		panic(fmt.Sprintf("Unsupported modifier type: %d", modType))
	}

	return &ModifierNode{
		BaseASTNode: NewBaseASTNode(),
		Type:        modType,
		Opcode:      opcode,
		Repeatable:  modType == AnnotationType,
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
	return mn.Type == GroovyParserPUBLIC || mn.Type == GroovyParserPROTECTED || mn.Type == GroovyParserPRIVATE
}

func (mn *ModifierNode) IsNonVisibilityModifier() bool {
	return mn.IsModifier() && !mn.IsVisibilityModifier()
}

func (mn *ModifierNode) IsAnnotation() bool {
	return mn.Type == AnnotationType
}

func (mn *ModifierNode) IsDef() bool {
	return mn.Type == GroovyParserDEF || mn.Type == GroovyParserVAR
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
