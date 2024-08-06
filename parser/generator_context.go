package parser

import (
	"strconv"
	"strings"
)

// GeneratorContext represents a context shared across generations of a class and its inner classes.
type GeneratorContext struct {
	innerClassIdx      int
	closureClassIdx    int
	syntheticMethodIdx int
	compileUnit        *CompileUnit
}

// NewGeneratorContext creates a new GeneratorContext instance.
func NewGeneratorContext(compileUnit *CompileUnit) *GeneratorContext {
	return &GeneratorContext{
		innerClassIdx:      1,
		closureClassIdx:    1,
		syntheticMethodIdx: 0,
		compileUnit:        compileUnit,
	}
}

// NewGeneratorContextWithOffset creates a new GeneratorContext instance with a specified inner class offset.
func NewGeneratorContextWithOffset(compileUnit *CompileUnit, innerClassOffset int) *GeneratorContext {
	return &GeneratorContext{
		innerClassIdx:      innerClassOffset,
		closureClassIdx:    1,
		syntheticMethodIdx: 0,
		compileUnit:        compileUnit,
	}
}

// GetNextInnerClassIdx returns and increments the inner class index.
func (gc *GeneratorContext) GetNextInnerClassIdx() int {
	idx := gc.innerClassIdx
	gc.innerClassIdx++
	return idx
}

// GetCompileUnit returns the compile unit.
func (gc *GeneratorContext) GetCompileUnit() *CompileUnit {
	return gc.compileUnit
}

// GetNextClosureInnerName generates the next closure inner name.
func (gc *GeneratorContext) GetNextClosureInnerName(owner, enclosingClass *ClassNode, enclosingMethod *MethodNode) string {
	return gc.getNextInnerName(owner, enclosingClass, enclosingMethod, "closure")
}

// GetNextLambdaInnerName generates the next lambda inner name.
func (gc *GeneratorContext) GetNextLambdaInnerName(owner, enclosingClass *ClassNode, enclosingMethod *MethodNode) string {
	return gc.getNextInnerName(owner, enclosingClass, enclosingMethod, "lambda")
}

func (gc *GeneratorContext) getNextInnerName(owner, enclosingClass *ClassNode, enclosingMethod *MethodNode, classifier string) string {
	methodName := ""
	if enclosingMethod != nil {
		methodName = enclosingMethod.name

		if enclosingClass.IsDerivedFrom(CLOSURE_TYPE) {
			methodName = ""
		} else {
			methodName = "_" + EncodeAsValidClassName(methodName)
		}
	}

	gc.closureClassIdx++
	return methodName + "_" + classifier + strconv.Itoa(gc.closureClassIdx-1)
}

// GetNextConstructorReferenceSyntheticMethodName generates the next constructor reference synthetic method name.
func (gc *GeneratorContext) GetNextConstructorReferenceSyntheticMethodName(enclosingMethodNode *MethodNode) string {
	gc.syntheticMethodIdx++
	if enclosingMethodNode == nil {
		return "ctorRef$" + strconv.Itoa(gc.syntheticMethodIdx-1)
	}
	return "ctorRef$" + strings.NewReplacer("<", "", ">", "").Replace(enclosingMethodNode.name) + "$" + strconv.Itoa(gc.syntheticMethodIdx-1)
}

var charactersToEncode = map[rune]bool{
	' ': true, '!': true, '"': true, '#': true, '$': true, '&': true, '\'': true,
	'(': true, ')': true, '*': true, '+': true, ',': true, '-': true, '.': true,
	'/': true, ':': true, ';': true, '<': true, '=': true, '>': true, '@': true,
	'[': true, '\\': true, ']': true, '^': true, '{': true, '}': true, '~': true,
}

// EncodeAsValidClassName encodes a string as a valid class name.
func EncodeAsValidClassName(name string) string {
	if name == "module-info" || name == "package-info" {
		return name
	}

	var b strings.Builder
	for _, ch := range name {
		if charactersToEncode[ch] {
			b.WriteRune('_')
		} else {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
