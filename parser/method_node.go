package parser

import (
	"fmt"
	"strings"
)

type MethodNode struct {
	AnnotatedNode
	name              string
	modifiers         int
	syntheticPublic   bool
	returnType        *ClassNode
	parameters        []*Parameter
	hasDefaultValue   bool
	code              Statement
	dynamicReturnType bool
	variableScope     *VariableScope
	exceptions        []*ClassNode
	genericsTypes     []*GenericsType
	typeDescriptor    string
}

func NewMethodNode(name string, modifiers int, returnType *ClassNode, parameters []*Parameter, exceptions []*ClassNode, code Statement) *MethodNode {
	mn := &MethodNode{
		name:       name,
		modifiers:  modifiers,
		exceptions: exceptions,
		code:       code,
	}
	mn.SetReturnType(returnType)
	mn.SetParameters(parameters)
	return mn
}

func (mn *MethodNode) GetTypeDescriptor() string {
	if mn.typeDescriptor == "" {
		mn.typeDescriptor = methodDescriptor(mn, false)
	}
	return mn.typeDescriptor
}

func (mn *MethodNode) invalidateCachedData() {
	mn.typeDescriptor = ""
}

// Getter and setter methods...

func (mn *MethodNode) IsAbstract() bool {
	return (mn.modifiers & ACC_ABSTRACT) != 0
}

func (mn *MethodNode) IsDefault() bool {
	return (mn.modifiers&(ACC_ABSTRACT|ACC_PUBLIC|ACC_STATIC) == ACC_PUBLIC) &&
		mn.GetDeclaringClass() != nil && mn.GetDeclaringClass().IsInterface()
}

func (mn *MethodNode) IsFinal() bool {
	return (mn.modifiers & ACC_FINAL) != 0
}

func (mn *MethodNode) IsStatic() bool {
	return (mn.modifiers & ACC_STATIC) != 0
}

func (mn *MethodNode) IsPublic() bool {
	return (mn.modifiers & ACC_PUBLIC) != 0
}

func (mn *MethodNode) IsPrivate() bool {
	return (mn.modifiers & ACC_PRIVATE) != 0
}

func (mn *MethodNode) IsProtected() bool {
	return (mn.modifiers & ACC_PROTECTED) != 0
}

func (mn *MethodNode) IsPackageScope() bool {
	return (mn.modifiers & (ACC_PUBLIC | ACC_PRIVATE | ACC_PROTECTED)) == 0
}

func (mn *MethodNode) GetFirstStatement() Statement {
	if mn.code == nil {
		return nil
	}
	first := mn.code
	for {
		if bs, ok := first.(*BlockStatement); ok {
			if len(bs.Statements) == 0 {
				return nil
			}
			first = bs.Statements[0]
		} else {
			break
		}
	}
	return first
}

func (mn *MethodNode) HasAnnotationDefault() bool {
	return mn.GetNodeMetaData("org.codehaus.groovy.ast.MethodNode.hasDefaultValue") == true
}

func (mn *MethodNode) SetAnnotationDefault(hasDefaultValue bool) {
	if hasDefaultValue {
		mn.PutNodeMetaData("org.codehaus.groovy.ast.MethodNode.hasDefaultValue", true)
	} else {
		mn.RemoveNodeMetaData("org.codehaus.groovy.ast.MethodNode.hasDefaultValue")
	}
}

func (mn *MethodNode) IsScriptBody() bool {
	return mn.GetNodeMetaData("org.codehaus.groovy.ast.MethodNode.isScriptBody") == true
}

func (mn *MethodNode) SetIsScriptBody() {
	mn.SetNodeMetaData("org.codehaus.groovy.ast.MethodNode.isScriptBody", true)
}

func (mn *MethodNode) IsStaticConstructor() bool {
	return mn.name == "<clinit>"
}

func (mn *MethodNode) IsConstructor() bool {
	return mn.name == "<init>"
}

func (mn *MethodNode) GetText() string {
	var mask int
	if _, ok := mn.(*ConstructorNode); ok {
		mask = ConstructorModifiers()
	} else {
		mask = MethodModifiers()
	}
	name := mn.GetName()
	if strings.Contains(name, " ") {
		name = fmt.Sprintf("\"%s\"", name)
	}
	return fmt.Sprintf("%s %s%s %s(%s)%s { ... }",
		GetModifiersText(mn.GetModifiers()&mask),
		ToGenericTypesString(mn.GetGenericsTypes()),
		GetClassText(mn.GetReturnType()),
		name,
		GetParametersText(mn.GetParameters()),
		GetThrowsClauseText(mn.GetExceptions()))
}

func (mn *MethodNode) String() string {
	declaringClass := mn.GetDeclaringClass()
	declaringClassStr := ""
	if declaringClass != nil {
		declaringClassStr = " from " + FormatTypeName(declaringClass)
	}
	return fmt.Sprintf("%s[%s%s]", mn.AnnotatedNode.String(), methodDescriptor(mn, true), declaringClassStr)
}

// Helper functions like methodDescriptor, GetModifiersText, etc. need to be implemented
