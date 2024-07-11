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
		mn.typeDescriptor = MethodDescriptor(mn, false)
	}
	return mn.typeDescriptor
}

func (mn *MethodNode) invalidateCachedData() {
	mn.typeDescriptor = ""
}

// Getter and setter methods...

func (mn *MethodNode) SetReturnType(returnType *ClassNode) {
	mn.invalidateCachedData()
	mn.dynamicReturnType = mn.dynamicReturnType || IsDynamicTyped(returnType)
	if returnType != nil {
		mn.returnType = returnType
	} else {
		mn.returnType = OBJECT_TYPE
	}
}

func (mn *MethodNode) SetParameters(parameters []*Parameter) {
	mn.invalidateCachedData()
	scope := NewVariableScope()
	mn.hasDefaultValue = false
	mn.parameters = parameters
	if parameters != nil && len(parameters) > 0 {
		for _, para := range parameters {
			if para.HasInitialExpression() {
				mn.hasDefaultValue = true
			}
			para.SetInStaticContext(mn.IsStatic())
			var varptr *Variable = &para.Variable
			scope.PutDeclaredVariable(varptr)
		}
	}
	mn.SetVariableScope(scope)
}

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
			if len(bs.statements) == 0 {
				return nil
			}
			first = bs.statements[0]
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

func ToGenericTypesString(genericsTypes []*GenericsType) string {
	if genericsTypes == nil || len(genericsTypes) == 0 {
		return ""
	}
	var parts []string
	for _, genericsType := range genericsTypes {
		parts = append(parts, genericsType.String())
	}
	return fmt.Sprintf("<%s> ", strings.Join(parts, ","))
}

func (mn *MethodNode) GetText() string {
	name := mn.name
	if strings.Contains(name, " ") {
		name = fmt.Sprintf("\"%s\"", name)
	}
	return fmt.Sprintf("%s %s%s %s(%s)%s { ... }",
		GetModifiersText(mn.modifiers),
		ToGenericTypesString(mn.genericsTypes),
		GetClassText(mn.returnType),
		name,
		GetParametersText(mn.parameters),
		GetThrowsClauseText(mn.exceptions))
}

func (mn *MethodNode) String() string {
	declaringClass := mn.GetDeclaringClass()
	declaringClassStr := ""
	if declaringClass != nil {
		declaringClassStr = " from " + declaringClass.GetText()
	}
	return fmt.Sprintf("%s[%s%s]", mn.AnnotatedNode.GetText(), MethodDescriptor(mn, true), declaringClassStr)
}

func (mn *MethodNode) SetVariableScope(variableScope *VariableScope) {
	mn.variableScope = variableScope
	variableScope.SetInStaticContext(mn.IsStatic())
}
