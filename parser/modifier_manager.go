package parser

import (
	"fmt"
	"reflect"
	"slices"
)

const (
	ABSTRACT = iota
	FINAL
	NATIVE
	STATIC
	VOLATILE
	SEALED
	NON_SEALED
)

type ModifierManager struct {
	astBuilder       *ASTBuilder
	modifierNodeList []*ModifierNode
}

var INVALID_MODIFIERS_MAP = map[reflect.Type][]int{
	reflect.TypeOf(&ConstructorNode{}): {STATIC, FINAL, ABSTRACT, NATIVE},
	reflect.TypeOf(&MethodNode{}):      {VOLATILE},
}

func NewModifierManager(astBuilder *ASTBuilder, modifierNodeList []*ModifierNode) *ModifierManager {
	mm := &ModifierManager{
		astBuilder: astBuilder,
	}
	mm.validate(modifierNodeList)
	mm.modifierNodeList = modifierNodeList
	return mm
}

func (mm *ModifierManager) AttachAnnotations(node *AnnotatedNode) *AnnotatedNode {
	for _, annotation := range mm.GetAnnotations() {
		node.AddAnnotationNode(annotation)
	}
	return node
}

func (mm *ModifierManager) ProcessVariableExpression(ve *VariableExpression) *VariableExpression {
	for _, e := range mm.modifierNodeList {
		ve.SetModifiers(ve.GetModifiers() | e.GetOpcode())
		// local variable does not attach annotations
	}
	return ve
}

func (mm *ModifierManager) ValidateConstructor(constructorNode *ConstructorNode) {
	mm.validateNode(INVALID_MODIFIERS_MAP[reflect.TypeOf(&ConstructorNode{})], astNodeAdapter{constructorNode})
}

func (mm *ModifierManager) ProcessMethodNode(mn *MethodNode) *MethodNode {
	for _, e := range mm.modifierNodeList {
		if e.IsVisibilityModifier() {
			mn.SetModifiers(mm.ClearVisibilityModifiers(mn.GetModifiers()) | e.GetOpcode())
		} else {
			mn.SetModifiers(mn.GetModifiers() | e.GetOpcode())
		}

		if e.IsAnnotation() {
			mn.AddAnnotation(e.GetAnnotationNode().GetClassNode())
		}
	}

	return mn
}

func (mm *ModifierManager) GetModifierCount() int {
	return len(mm.modifierNodeList)
}

func (mm *ModifierManager) validate(modifierNodeList []*ModifierNode) {
	modifierNodeCounter := make(map[*ModifierNode]int)
	visibilityModifierCnt := 0

	for _, modifierNode := range modifierNodeList {
		cnt, exists := modifierNodeCounter[modifierNode]
		if !exists {
			modifierNodeCounter[modifierNode] = 1
		} else if cnt == 1 && !modifierNode.IsRepeatable() {
			panic(createParsingFailedException(fmt.Sprintf("Cannot repeat modifier[%s]", modifierNode.GetText()), astNodeAdapter{modifierNode}))
		}

		if modifierNode.IsVisibilityModifier() {
			visibilityModifierCnt++
			if visibilityModifierCnt > 1 {
				panic(createParsingFailedException(fmt.Sprintf("Cannot specify modifier[%s] when access scope has already been defined", modifierNode.GetText()), astNodeAdapter{modifierNode}))
			}
		}
	}
}

func (mm *ModifierManager) Validate(node interface{}) {
	switch n := node.(type) {
	case *MethodNode:
		mm.validateMethod(n)
	case *ConstructorNode:
		mm.validateConstructor(n)
	}
}

func (mm *ModifierManager) validateMethod(methodNode *MethodNode) {
	mm.validateNode(INVALID_MODIFIERS_MAP[reflect.TypeOf(&MethodNode{})], astNodeAdapter{methodNode})
}

func (mm *ModifierManager) validateConstructor(constructorNode *ConstructorNode) {
	mm.validateNode(INVALID_MODIFIERS_MAP[reflect.TypeOf(&ConstructorNode{})], astNodeAdapter{constructorNode})
}

func (mm *ModifierManager) validateNode(invalidModifierList []int, node SourcePosition) {
	for _, e := range mm.modifierNodeList {
		if slices.Contains(invalidModifierList, e.GetType()) {
			panic(createParsingFailedException(fmt.Sprintf("%s has an incorrect modifier '%s'.", reflect.TypeOf(node).Name(), e), node))
		}
	}
}

func (mm *ModifierManager) calcModifiersOpValue(t int) int {
	result := 0
	for _, modifierNode := range mm.modifierNodeList {
		result |= modifierNode.GetOpcode()
	}

	if !mm.ContainsVisibilityModifier() {
		if t == 1 {
			result |= ACC_SYNTHETIC | ACC_PUBLIC
		} else if t == 2 {
			result |= ACC_PUBLIC
		}
	}

	return result
}

func (mm *ModifierManager) GetClassModifiersOpValue() int {
	return mm.calcModifiersOpValue(1)
}

func (mm *ModifierManager) GetClassMemberModifiersOpValue() int {
	return mm.calcModifiersOpValue(2)
}

func (mm *ModifierManager) GetAnnotations() []*AnnotationNode {
	var annotations []*AnnotationNode
	for _, m := range mm.modifierNodeList {
		if m.IsAnnotation() {
			annotations = append(annotations, m.GetAnnotationNode())
		}
	}
	return annotations
}

func (mm *ModifierManager) ContainsAny(modifierTypes ...int) bool {
	for _, e := range mm.modifierNodeList {
		for _, modifierType := range modifierTypes {
			if modifierType == e.GetType() {
				return true
			}
		}
	}
	return false
}

func (mm *ModifierManager) Get(modifierType int) *ModifierNode {
	for _, e := range mm.modifierNodeList {
		if modifierType == e.GetType() {
			return e
		}
	}
	return nil
}

func (mm *ModifierManager) ContainsAnnotations() bool {
	for _, m := range mm.modifierNodeList {
		if m.IsAnnotation() {
			return true
		}
	}
	return false
}

func (mm *ModifierManager) ContainsVisibilityModifier() bool {
	for _, m := range mm.modifierNodeList {
		if m.IsVisibilityModifier() {
			return true
		}
	}
	return false
}

func (mm *ModifierManager) ContainsNonVisibilityModifier() bool {
	for _, m := range mm.modifierNodeList {
		if m.IsNonVisibilityModifier() {
			return true
		}
	}
	return false
}

func (mm *ModifierManager) ProcessParameter(parameter *Parameter) *Parameter {
	for _, e := range mm.modifierNodeList {
		parameter.SetModifiers(parameter.GetModifiers() | e.GetOpcode())
		if e.IsAnnotation() {
			parameter.AddAnnotationNode(e.GetAnnotationNode())
		}
	}
	return parameter
}

func (mm *ModifierManager) ClearVisibilityModifiers(modifiers int) int {
	return modifiers & ^ACC_PUBLIC & ^ACC_PROTECTED & ^ACC_PRIVATE
}
