package parser

import (
	"fmt"
	"strings"
)

var _ IClassNode = (*ClassNode)(nil)

var EMPTY_ARRAY = []IClassNode{}
var THIS = OBJECT_TYPE
var SUPER = OBJECT_TYPE

type IClassNode interface {
	NodeMetaDataHandler
	ASTNode
	IsInterface() bool
	GetModifiers() int
	GetEnclosingMethod() MethodOrConstructorNode
	GetOuterClass() IClassNode
	GetGenericsTypes() []*GenericsType
	AsGenericsType() *GenericsType
	IsGenericsPlaceHolder() bool
	IsArray() bool
	GetComponentType() *ClassNode
	GetName() string
	GetUnresolvedName() string
	IsStatic() bool
	RenameField(oldName, newName string)
	Equals(other IClassNode) bool
	Redirect() *ClassNode
	MakeArray() *ClassNode
	IsUsingGenerics() bool
	SetGenericsPlaceHolder(placeholder bool)
	IsDerivedFrom(type_ IClassNode) bool
	ImplementsInterface(classNode IClassNode) bool
	GetUnresolvedInterfaces() []IClassNode
	GetUnresolvedSuperClass() IClassNode
	GetPlainNodeReference() IClassNode
	SetGenericsTypes(genericsTypes []*GenericsType)
	GetSuperClass() IClassNode
	DeclaresInterface(classNode IClassNode) bool
	IsRedirectNode() bool
	ToStringWithRedirect(showRedirect bool) string
	GetDeclaredMethods(name string) []MethodOrConstructorNode
	GetAllInterfaces() map[IClassNode]bool
	GetAbstractMethods() []MethodOrConstructorNode
	GetDeclaredMethodsMap() map[string]MethodOrConstructorNode
	GetAllDeclaredMethods() []MethodOrConstructorNode
	GetAnnotations() []*AnnotationNode
	GetDeclaredField(name string) *FieldNode
	GetInterfaces() []IClassNode
	GetAllInterfacesWithSet(set map[IClassNode]bool)
	GetAnnotationsWithType(type_ IClassNode) []*AnnotationNode
	GetDeclaredConstructors() []MethodOrConstructorNode
	GetMethods(name string) []MethodOrConstructorNode
	IsResolved() bool
	SetRedirect(node IClassNode)
	GetNameWithoutPackage() string
	SetScript(isScript bool)
	SetScriptBody(isScriptBody bool)
	SetSuperClass(superClass IClassNode)
	GetUnresolvedInterfacesWithDeref(deref bool) []IClassNode
	GetUnresolvedSuperClassWithDeref(deref bool) IClassNode
	GetMixins() []*MixinNode
	SetMixins(mixins []*MixinNode)
	SetName(name string) string
	GetProperties() []*PropertyNode
	GetFields() []*FieldNode
	SetInterfaces(interfaces []IClassNode)
	SetEnclosingMethod(enclosingMethod MethodOrConstructorNode)
	AddInterface(iface IClassNode)
	AddInnerClass(innerClass *InnerClassNode)
	AddField(node *FieldNode)
	SetSyntheticPublic(isSyntheticPublic bool)
	AddAnnotations(annotations []*AnnotationNode)
	AddAnnotationNode(annotation *AnnotationNode)
	SetModifiers(modifiers int)
	GetInnerClasses() []*InnerClassNode
	IsEnum() bool
	SetUsingGenerics(usesGenerics bool)
	AddStaticInitializerStatements(statements []Statement, fieldInit bool)
	AddObjectInitializerStatements(statement Statement)
	AddTypeAnnotations(annotations []*AnnotationNode)
	GetAnnotationsOfType(typ IClassNode) []*AnnotationNode
	AddAnnotation(typ IClassNode) *AnnotationNode
	IsAbstract() bool
	IsAnnotationDefinition() bool
	AddMethod(node MethodOrConstructorNode)
	AddConstructorWithDetails(modifiers int, parameters []*Parameter, exceptions []IClassNode, code Statement) *ConstructorNode
	HasProperty(name string) bool
	RemoveField(fieldToRemove *FieldNode)
	AddProperty(node *PropertyNode)
	GetProperty(name string) *PropertyNode
	GetPlainNodeReferenceHelper(skipPrimitives bool) *ClassNode
	GetGetterMethod(getterName string, searchSupers bool) MethodOrConstructorNode
	GetDeclaredMethod(name string, parameters []*Parameter) MethodOrConstructorNode
	IsInnerClass() bool
	GetOuterClasses() []IClassNode
	GetPackageName() string
	HasPossibleMethod(name string, arguments Expression) bool
	AddConstructor(node *ConstructorNode)
}

type ClassNode struct {
	*AnnotatedNode
	name                string
	modifiers           int
	superClass          IClassNode
	interfaces          []IClassNode
	fields              []*FieldNode
	methods             map[string][]MethodOrConstructorNode
	constructors        []MethodOrConstructorNode
	isPrimaryNode       bool
	redirect            IClassNode
	componentType       *ClassNode
	genericsTypes       []*GenericsType
	usingGenerics       bool
	permittedSubclasses []*ClassNode
	recordComponents    []*RecordComponentNode
	placeholder         bool
	script              bool
	scriptBody          bool
	fieldIndex          map[string]*FieldNode
	innerClasses        []*InnerClassNode
	syntheticPublic     bool
	enclosingMethod     MethodOrConstructorNode
	typeAnnotations     []*AnnotationNode
	annotated           bool
	objectInitializers  []Statement
	properties          []*PropertyNode
	mixins              []*MixinNode
}

func NewClassNode(name string, modifiers int, superClass IClassNode) *ClassNode {
	return &ClassNode{
		AnnotatedNode: NewAnnotatedNode(),
		name:          name,
		modifiers:     modifiers,
		superClass:    superClass,
		methods:       make(map[string][]MethodOrConstructorNode),
		fieldIndex:    make(map[string]*FieldNode),
	}
}

func NewClassNodeWithInterfaces(name string, modifiers int, superClass IClassNode, interfaces []IClassNode, mixins []*MixinNode) *ClassNode {
	cn := NewClassNode(name, modifiers, superClass)
	cn.interfaces = interfaces
	// Assuming you want to add mixins to the ClassNode, you might need to add a mixins field to the ClassNode struct
	// cn.mixins = mixins
	return cn
}

func (cn *ClassNode) GetPackageName() string {
	name := cn.GetName()
	idx := strings.LastIndex(name, ".")
	if idx > 0 {
		return name[:idx]
	}
	return ""
}

func (cn *ClassNode) GetGetterMethod(getterName string, searchSupers bool) MethodOrConstructorNode {
	var getterMethod MethodOrConstructorNode

	isNullOrSynthetic := func(method MethodOrConstructorNode) bool {
		return method == nil || (method.GetModifiers()&ACC_SYNTHETIC) != 0
	}

	booleanReturnOnly := strings.HasPrefix(getterName, "is")
	for _, method := range cn.GetDeclaredMethods(getterName) {
		if booleanReturnOnly && !IsPrimitiveBoolean(method.GetReturnType()) {
			continue
		}
		if !booleanReturnOnly && method.IsVoidMethod() {
			continue
		}
		if len(method.GetParameters()) == 0 {
			if isNullOrSynthetic(getterMethod) {
				getterMethod = method
			}
		} else if method.HasDefaultValue() && AllParametersHaveInitialExpression(method.GetParameters()) {
			if isNullOrSynthetic(getterMethod) {
				getterMethod = NewMethodNode(method.GetName(), method.GetModifiers()&^ACC_ABSTRACT, method.GetReturnType(), []*Parameter{}, method.GetExceptions(), nil)
				getterMethod.SetSynthetic(true)
				getterMethod.SetDeclaringClass(cn)
				getterMethod.AddAnnotations(method.GetAnnotations())
				MarkAsGenerated(cn, getterMethod)
				getterMethod.SetGenericsTypes(method.GetGenericsTypes())
			}
		}
	}

	if searchSupers && isNullOrSynthetic(getterMethod) {
		superClass := cn.GetSuperClass()
		if superClass != nil {
			method := superClass.GetGetterMethod(getterName, true)
			if getterMethod == nil || !isNullOrSynthetic(method) {
				getterMethod = method
			}
		}

		if getterMethod == nil && len(cn.GetInterfaces()) > 0 {
			for anInterface := range cn.GetAllInterfaces() {
				method := anInterface.GetDeclaredMethod(getterName, []*Parameter{})
				if method != nil && method.IsDefault() {
					if booleanReturnOnly && IsPrimitiveBoolean(method.GetReturnType()) {
						getterMethod = method
						break
					}
					if !booleanReturnOnly && !method.IsVoidMethod() {
						getterMethod = method
						break
					}
				}
			}
		}
	}

	return getterMethod
}

func AllParametersHaveInitialExpression(parameters []*Parameter) bool {
	for _, param := range parameters {
		if !param.HasInitialExpression() {
			return false
		}
	}
	return true
}

func MarkAsGenerated(cn *ClassNode, mn MethodOrConstructorNode) {
	// TODO: implement
}

func (cn *ClassNode) GetMethods(name string) []MethodOrConstructorNode {
	var list []MethodOrConstructorNode
	var node IClassNode = cn

	for node != nil {
		list = append(list, node.GetDeclaredMethods(name)...)
		node = node.GetSuperClass()
	}

	return list
}

func (cn *ClassNode) GetDeclaredConstructors() []MethodOrConstructorNode {
	if cn.redirect != nil {
		return cn.redirect.GetDeclaredConstructors()
	}
	cn.lazyClassInit()
	if cn.constructors == nil {
		cn.constructors = make([]MethodOrConstructorNode, 0)
	}
	return cn.constructors
}

func (cn *ClassNode) GetDeclaredConstructor(parameters []*Parameter) MethodOrConstructorNode {
	for _, constructor := range cn.GetDeclaredConstructors() {
		if ParametersEqual(constructor.GetParameters(), parameters) {
			return constructor
		}
	}
	return nil
}

func (cn *ClassNode) GetUnresolvedInterfaces() []IClassNode {
	return cn.GetUnresolvedInterfacesWithDeref(true)
}

func (cn *ClassNode) GetUnresolvedInterfacesWithDeref(deref bool) []IClassNode {
	if deref {
		if cn.redirect != nil {
			return cn.redirect.GetUnresolvedInterfacesWithDeref(true)
		}
		cn.lazyClassInit()
	}
	return cn.interfaces
}

func (cn *ClassNode) GetUnresolvedSuperClass() IClassNode {
	return cn.GetUnresolvedSuperClassWithDeref(true)
}

func (cn *ClassNode) GetUnresolvedSuperClassWithDeref(deref bool) IClassNode {
	if deref {
		if cn.redirect != nil {
			return cn.redirect.GetUnresolvedSuperClassWithDeref(true)
		}
		cn.lazyClassInit()
	}
	return cn.superClass
}

func (cn *ClassNode) GetMixins() []*MixinNode {
	if cn.redirect != nil {
		return cn.redirect.GetMixins()
	}
	return cn.mixins
}

func (cn *ClassNode) SetMixins(mixins []*MixinNode) {
	if cn.redirect != nil {
		cn.redirect.SetMixins(mixins)
	} else {
		cn.mixins = mixins
	}
}

func (cn *ClassNode) SetName(name string) string {
	if cn.redirect != nil {
		return cn.redirect.SetName(name)
	}
	cn.name = name
	return name
}

func (cn *ClassNode) GetProperty(name string) *PropertyNode {
	for _, prop := range cn.GetProperties() {
		if prop.GetName() == name {
			return prop
		}
	}
	return nil
}

func (cn *ClassNode) AddProperty(node *PropertyNode) {
	node.SetDeclaringClass(cn.Redirect())
	cn.AddField(node.GetField())
	properties := cn.GetProperties()
	cn.properties = append(properties, node)
}

// Add this method to the ClassNode struct
func (cn *ClassNode) RemoveField(fieldToRemove *FieldNode) {
	r := cn.Redirect()
	for i, field := range r.fields {
		if field == fieldToRemove {
			r.fields = append(r.fields[:i], r.fields[i+1:]...)
			delete(r.fieldIndex, fieldToRemove.GetName())
			return
		}
	}
}

func (cn *ClassNode) HasProperty(name string) bool {
	for _, prop := range cn.GetProperties() {
		if prop.GetName() == name {
			return true
		}
	}
	return false
}

func (cn *ClassNode) GetProperties() []*PropertyNode {
	if cn.redirect != nil {
		return cn.redirect.GetProperties()
	}
	cn.lazyClassInit()
	if cn.properties == nil {
		cn.properties = make([]*PropertyNode, 0)
	}
	return cn.properties
}

func (cn *ClassNode) SetModifiers(modifiers int) {
	cn.modifiers = modifiers
}

func (cn *ClassNode) GetFields() []*FieldNode {
	if cn.redirect != nil {
		return cn.redirect.GetFields()
	}
	cn.lazyClassInit()
	if cn.fields == nil {
		cn.fields = make([]*FieldNode, 0)
	}
	return cn.fields
}

func (cn *ClassNode) SetInterfaces(interfaces []IClassNode) {
	if cn.redirect != nil {
		cn.redirect.SetInterfaces(interfaces)
	} else {
		cn.interfaces = interfaces
		// GROOVY-10763: update generics indicator
		if interfaces != nil && !cn.usingGenerics && cn.isPrimaryNode {
			for _, iface := range interfaces {
				cn.usingGenerics = cn.usingGenerics || iface.IsUsingGenerics()
			}
		}
	}
}

func (cn *ClassNode) GetOrAddStaticInitializer() MethodOrConstructorNode {
	const classInitializer = "<clinit>"
	declaredMethods := cn.GetDeclaredMethods(classInitializer)

	if len(declaredMethods) == 0 {
		method := cn.AddMethodWithDetails(
			classInitializer,
			ACC_STATIC,
			VOID_TYPE,
			[]*Parameter{},
			EMPTY_ARRAY,
			NewBlockStatement(),
		)
		method.SetSynthetic(true)
		return method
	}

	return declaredMethods[0]
}

func (cn *ClassNode) AddObjectInitializerStatements(statement Statement) {
	statements := cn.GetObjectInitializerStatements()
	cn.objectInitializers = append(statements, statement)
}

func (cn *ClassNode) AddStaticInitializerStatements(statements []Statement, fieldInit bool) {
	method := cn.GetOrAddStaticInitializer()
	block := GetCodeAsBlock(method)

	if !fieldInit {
		block.AddStatements(statements)
	} else {
		blockStatements := block.GetStatements()
		statements = append(statements, blockStatements...)
		block.ClearStatements()
		block.AddStatements(statements)
	}
}

func (cn *ClassNode) GetObjectInitializerStatements() []Statement {
	r := cn.Redirect()
	if r.objectInitializers == nil {
		r.objectInitializers = make([]Statement, 0)
	}
	return r.objectInitializers
}

func (cn *ClassNode) IsAnnotated() bool {
	return cn.annotated
}

func (cn *ClassNode) SetAnnotated(annotated bool) {
	cn.annotated = annotated
}

func (cn *ClassNode) IsRedirectNode() bool {
	return cn.redirect != nil
}

func (cn *ClassNode) IsResolved() bool {
	// Implement this method based on your specific requirements
	// For now, we'll return false as a placeholder
	return false
}

func (cn *ClassNode) IsPrimaryClassNode() bool {
	return cn.isPrimaryNode
}

func (cn *ClassNode) AddTypeAnnotation(annotation *AnnotationNode) {
	if !cn.IsRedirectNode() && (cn.IsResolved() || cn.IsPrimaryClassNode()) {
		panic(fmt.Sprintf("Adding type annotation @%s to non-redirect node: %s", annotation.GetClassNode().GetNameWithoutPackage(), cn.GetName()))
	}
	if cn.typeAnnotations == nil {
		cn.typeAnnotations = make([]*AnnotationNode, 0, 3)
	}
	cn.typeAnnotations = append(cn.typeAnnotations, annotation)
	cn.SetAnnotated(true)
}

func (cn *ClassNode) AddTypeAnnotations(annotations []*AnnotationNode) {
	for _, annotation := range annotations {
		cn.AddTypeAnnotation(annotation)
	}
}

func (cn *ClassNode) GetEnclosingMethod() MethodOrConstructorNode {
	if cn.redirect != nil {
		return cn.redirect.GetEnclosingMethod()
	}
	return cn.enclosingMethod
}

func (cn *ClassNode) SetEnclosingMethod(enclosingMethod MethodOrConstructorNode) {
	if cn.redirect != nil {
		cn.redirect.SetEnclosingMethod(enclosingMethod)
	} else {
		cn.enclosingMethod = enclosingMethod
	}
}

func (cn *ClassNode) AddConstructor(node *ConstructorNode) {
	r := cn.Redirect()
	node.SetDeclaringClass(r)
	if r.constructors == nil {
		r.constructors = make([]MethodOrConstructorNode, 0, 4)
	}
	r.constructors = append(r.constructors, node)
}

func (cn *ClassNode) AddConstructorWithDetails(modifiers int, parameters []*Parameter, exceptions []IClassNode, code Statement) *ConstructorNode {
	node := NewConstructorNodeWithParams(modifiers, parameters, exceptions, code)
	cn.AddConstructor(node)
	return node
}

func (cn *ClassNode) GetDeclaredMethod(name string, parameters []*Parameter) MethodOrConstructorNode {
	zeroParameters := len(parameters) == 0
	for _, method := range cn.GetDeclaredMethods(name) {
		if zeroParameters {
			if len(method.GetParameters()) == 0 {
				return method
			}
		} else {
			if ParametersEqual(method.GetParameters(), parameters) {
				return method
			}
		}
	}
	return nil
}

func (cn *ClassNode) GetDeclaredMethods(name string) []MethodOrConstructorNode {
	if cn.redirect != nil {
		return cn.redirect.GetDeclaredMethods(name)
	}
	return cn.methods[name]
}

func (cn *ClassNode) AddMethod(node MethodOrConstructorNode) {
	r := cn.Redirect()
	node.SetDeclaringClass(r)
	if r.methods == nil {
		r.methods = make(map[string][]MethodOrConstructorNode)
	}
	r.methods[node.GetName()] = append(r.methods[node.GetName()], node)
}

func (cn *ClassNode) AddMethodWithDetails(name string, modifiers int, returnType *ClassNode, parameters []*Parameter, exceptions []IClassNode, code Statement) MethodOrConstructorNode {
	other := cn.GetDeclaredMethod(name, parameters)
	// don't add duplicate methods
	if other != nil {
		return other
	}
	node := NewMethodNode(name, modifiers, returnType, parameters, exceptions, code)
	cn.AddMethod(node)
	return node
}

func (cn *ClassNode) SetSyntheticPublic(isSyntheticPublic bool) {
	cn.syntheticPublic = isSyntheticPublic
}

func (cn *ClassNode) IsInnerClass() bool {
	return false
}

func (cn *ClassNode) Visit(vistor GroovyCodeVisitor) {
	return
}

func (cn *ClassNode) IsDerivedFrom(type_ IClassNode) bool {
	if IsPrimitiveVoid(cn) {
		return IsPrimitiveVoid(type_)
	}
	if IsObjectType(type_) {
		return true
	}
	if cn.IsArray() && type_.IsArray() &&
		IsObjectType(type_.GetComponentType()) &&
		!IsPrimitiveType(cn.GetComponentType()) {
		return true
	}
	var node IClassNode
	for node = cn; node != nil; node = node.GetSuperClass() {
		if type_.Equals(node) {
			return true
		}
	}
	return false
}

func (cn *ClassNode) MakeArray() *ClassNode {
	if cn.redirect != nil {
		node := cn.redirect.MakeArray()
		node.componentType = cn
		return node
	}

	// Note: Go doesn't have direct equivalents for Java's reflection.
	// You might need to adjust this part based on your specific needs.
	node := NewClassNode(cn.name+"[]", cn.modifiers, cn.superClass)
	node.componentType = cn
	return node
}

func (cn *ClassNode) Equals(other IClassNode) bool {
	if cn == other {
		return true
	}
	if cn.redirect != nil {
		return cn.redirect.Equals(other)
	}
	if cn.componentType != nil {
		return cn.componentType.Equals(other.GetComponentType())
	}
	return other.GetName() == cn.GetName() // arrays could be "T[]" or "[LT;"
}

func (cn *ClassNode) GetText() string {
	return cn.GetName()
}

func (cn *ClassNode) GetName() string {
	if cn.redirect != nil {
		return cn.redirect.GetName()
	}
	return cn.name
}

func (cn *ClassNode) GetUnresolvedName() string {
	return cn.name
}

func (cn *ClassNode) GetNameWithoutPackage() string {
	name := cn.GetName()
	idx := strings.LastIndex(name, ".")
	if idx > 0 {
		return name[idx+1:]
	}
	return name
}

func (cn *ClassNode) IsPrimaryNode() bool {
	return cn.isPrimaryNode
}

func (cn *ClassNode) SetPrimaryNode(isPrimary bool) {
	cn.isPrimaryNode = isPrimary
}

func (cn *ClassNode) GetSuperClass() IClassNode {
	if cn.redirect != nil {
		return cn.redirect.GetSuperClass()
	}
	return cn.superClass
}

func (cn *ClassNode) SetSuperClass(superClass IClassNode) {
	if cn.redirect != nil {
		cn.redirect.SetSuperClass(superClass)
	} else {
		cn.superClass = superClass
	}
}

func (cn *ClassNode) GetInnerClasses() []*InnerClassNode {
	return cn.innerClasses
}

func (cn *ClassNode) GetInterfaces() []IClassNode {
	if cn.redirect != nil {
		return cn.redirect.GetInterfaces()
	}
	return cn.interfaces
}

func (cn *ClassNode) AddInterface(iface IClassNode) {
	if cn.redirect != nil {
		cn.redirect.AddInterface(iface)
	} else {
		cn.interfaces = append(cn.interfaces, iface)
	}
}

func (cn *ClassNode) IsInterface() bool {
	return (cn.modifiers & ACC_INTERFACE) != 0
}

func (cn *ClassNode) IsArray() bool {
	return cn.componentType != nil
}

func (cn *ClassNode) GetModifiers() int {
	return cn.modifiers
}

func (cn *ClassNode) IsEnum() bool {
	return (cn.modifiers & ACC_ENUM) != 0
}

func (cn *ClassNode) IsStatic() bool {
	return (cn.modifiers & ACC_STATIC) != 0
}

func (cn *ClassNode) IsAnnotationDefinition() bool {
	return cn.IsInterface() && (cn.modifiers&ACC_ANNOTATION) != 0
}

func (cn *ClassNode) IsAbstract() bool {
	return (cn.modifiers & ACC_ABSTRACT) != 0
}

func (cn *ClassNode) GetComponentType() *ClassNode {
	return cn.componentType
}

func (cn *ClassNode) IsGenericsPlaceHolder() bool {
	return cn.placeholder
}

func (cn *ClassNode) GetPlainNodeReferenceHelper(skipPrimitives bool) *ClassNode {
	if skipPrimitives && IsPrimitiveType(cn) {
		return cn
	}

	n := NewClassNode(cn.name, cn.modifiers, cn.superClass)
	n.isPrimaryNode = false
	n.SetRedirect(cn.Redirect())

	if cn.IsArray() {
		n.componentType = cn.redirect.GetComponentType()
	}

	return n
}

func (cn *ClassNode) SetRedirect(node IClassNode) {
	if cn.isPrimaryNode {
		panic("tried to set a redirect for a primary ClassNode (" + cn.GetName() + "->" + node.GetName() + ").")
	}
	if node != nil && !cn.IsGenericsPlaceHolder() {
		node = node.Redirect()
	}
	if node == cn {
		return
	}
	cn.redirect = node
}

func (cn *ClassNode) Redirect() *ClassNode {
	if cn.redirect != nil {
		return cn.redirect.Redirect()
	}
	return cn
}

func (cn *ClassNode) GetPlainNodeReference() IClassNode {
	return cn.GetPlainNodeReferenceHelper(true)
}

func (cn *ClassNode) SetScript(isScript bool) {
	if cn.redirect != nil {
		cn.redirect.SetScript(isScript)
	} else {
		cn.script = isScript
	}
}

func (cn *ClassNode) SetScriptBody(isScriptBody bool) {
	if cn.redirect != nil {
		cn.redirect.SetScriptBody(isScriptBody)
	} else {
		cn.scriptBody = isScriptBody
	}
}

func (cn *ClassNode) ImplementsInterface(classNode IClassNode) bool {
	var node IClassNode
	node = cn.Redirect()
	for node != nil {
		if node.DeclaresInterface(classNode) {
			return true
		}
		node = node.GetSuperClass()
	}
	return false
}

func (cn *ClassNode) DeclaresAnyInterfaces(classNodes ...IClassNode) bool {
	for _, classNode := range classNodes {
		if cn.DeclaresInterface(classNode) {
			return true
		}
	}
	return false
}

func (cn *ClassNode) DeclaresInterface(classNode IClassNode) bool {
	interfaces := cn.GetInterfaces()
	for _, face := range interfaces {
		if face.Equals(classNode) {
			return true
		}
	}
	for _, face := range interfaces {
		if face.DeclaresInterface(classNode) {
			return true
		}
	}
	return false
}

func (cn *ClassNode) GetDeclaredField(name string) *FieldNode {
	if cn.redirect != nil {
		return cn.redirect.GetDeclaredField(name)
	}
	cn.lazyClassInit()
	if cn.fieldIndex == nil {
		return nil
	}
	return cn.fieldIndex[name]
}

func (cn *ClassNode) RenameField(oldName, newName string) {
	r := cn.Redirect()
	if r.fieldIndex != nil {
		if field, exists := r.fieldIndex[oldName]; exists {
			delete(r.fieldIndex, oldName)
			r.fieldIndex[newName] = field
		}
	}
}

func (cn *ClassNode) AddField(node *FieldNode) {
	cn.addField(node, true)
}

func (cn *ClassNode) addField(node *FieldNode, doappend bool) {
	r := cn.Redirect()
	node.SetDeclaringClass(r)
	node.SetOwner(r)
	if r.fields == nil {
		r.fields = make([]*FieldNode, 0, 4)
	}
	if r.fieldIndex == nil {
		r.fieldIndex = make(map[string]*FieldNode)
	}

	if doappend {
		r.fields = append(r.fields, node)
	} else {
		r.fields = append([]*FieldNode{node}, r.fields...)
	}
	r.fieldIndex[node.GetName()] = node
}

func (cn *ClassNode) AddFieldWithInitialValue(name string, modifiers int, type_ *ClassNode, initialValue Expression) *FieldNode {
	node := NewFieldNode(name, modifiers, type_, cn.Redirect(), initialValue)
	cn.AddField(node)
	return node
}

func (cn *ClassNode) AddInnerClass(innerClass *InnerClassNode) {
	if cn.redirect != nil {
		cn.redirect.AddInnerClass(innerClass)
	} else {
		if cn.innerClasses == nil {
			cn.innerClasses = make([]*InnerClassNode, 0, 4)
		}
		cn.innerClasses = append(cn.innerClasses, innerClass)
	}
}

func (cn *ClassNode) lazyClassInit() {
	// Implement the lazy initialization logic here
	// This method should populate the fieldIndex map if it hasn't been done yet
	// For now, we'll leave it as an empty implementation
}

func (cn *ClassNode) GetGenericsTypes() []*GenericsType {
	return cn.genericsTypes
}

func (cn *ClassNode) AsGenericsType() *GenericsType {
	if !cn.IsGenericsPlaceHolder() {
		return NewGenericsTypeWithBasicType(cn)
	} else if cn.genericsTypes != nil && len(cn.genericsTypes) > 0 && cn.genericsTypes[0].GetUpperBounds() != nil {
		return cn.genericsTypes[0]
	} else {
		var upper IClassNode = cn
		if cn.redirect != nil {
			upper = cn.redirect
		}
		return NewGenericsType(cn, []IClassNode{upper}, nil)
	}
}

func (cn *ClassNode) SetGenericsPlaceHolder(placeholder bool) {
	if cn.redirect != nil {
		cn.redirect.SetGenericsPlaceHolder(placeholder)
	} else {
		cn.usingGenerics = cn.usingGenerics || placeholder
		cn.placeholder = placeholder
	}
}

func (cn *ClassNode) GetOuterClass() IClassNode {
	if cn.redirect != nil {
		return cn.redirect.GetOuterClass()
	}
	return nil
}

func (cn *ClassNode) GetOuterClasses() []IClassNode {
	outer := cn.GetOuterClass()
	if outer == nil {
		return []IClassNode{}
	}

	result := make([]IClassNode, 0, 4)
	for outer != nil {
		result = append(result, outer)
		outer = outer.GetOuterClass()
	}

	return result
}

func (cn *ClassNode) SetGenericsTypes(genericsTypes []*GenericsType) {
	if cn.redirect != nil {
		cn.redirect.SetGenericsTypes(genericsTypes)
	} else {
		cn.usingGenerics = cn.usingGenerics || genericsTypes != nil
		cn.genericsTypes = genericsTypes
	}
}

func (cn *ClassNode) IsUsingGenerics() bool {
	return cn.usingGenerics
}

func (cn *ClassNode) SetUsingGenerics(usesGenerics bool) {
	cn.usingGenerics = usesGenerics
}

func (cn *ClassNode) ToString() string {
	return cn.ToStringWithRedirect(true)
}

func (cn *ClassNode) ToStringWithRedirect(showRedirect bool) string {
	if cn.IsArray() {
		return cn.GetComponentType().ToStringWithRedirect(showRedirect) + "[]"
	}

	placeholder := cn.IsGenericsPlaceHolder()
	var ret strings.Builder

	if !placeholder {
		ret.WriteString(cn.GetName())
	} else {
		ret.WriteString(cn.GetUnresolvedName())
	}

	genericsTypes := cn.GetGenericsTypes()
	if !placeholder && genericsTypes != nil {
		ret.WriteString("<")
		for i, gt := range genericsTypes {
			if i != 0 {
				ret.WriteString(", ")
			}
			ret.WriteString(gt.String())
		}
		ret.WriteString(">")
	}

	if showRedirect && cn.redirect != nil {
		ret.WriteString(" -> ")
		ret.WriteString(cn.redirect.ToStringWithRedirect(showRedirect))
	}

	return ret.String()
}

func (cn *ClassNode) HasPossibleMethod(name string, arguments Expression) bool {
	var count int
	if tupleExpr, ok := arguments.(*TupleExpression); ok {
		count = len(tupleExpr.GetExpressions())
	} else {
		count = 0
	}

	var node IClassNode

	for node = cn; node != nil; node = node.GetSuperClass() {
		for _, mn := range node.GetDeclaredMethods(name) {
			if !mn.IsStatic() && HasCompatibleNumberOfArgs(mn, count) {
				return true
			}
		}
		for iface := range node.GetAllInterfaces() {
			for _, mn := range iface.GetDeclaredMethods(name) {
				if mn.IsDefault() && HasCompatibleNumberOfArgs(mn, count) {
					return true
				}
			}
		}
	}

	return false
}

func HasCompatibleNumberOfArgs(mn MethodOrConstructorNode, count int) bool {
	params := mn.GetParameters()
	lastParamIndex := len(params) - 1
	return count == len(params) || (count >= lastParamIndex && isPotentialVarArg(mn, lastParamIndex))
}

func isPotentialVarArg(mn MethodOrConstructorNode, lastParamIndex int) bool {
	return lastParamIndex >= 0 && mn.GetParameters()[lastParamIndex].GetType().IsArray()
}

func (cn *ClassNode) GetAllInterfaces() map[IClassNode]bool {
	result := make(map[IClassNode]bool)
	if cn.IsInterface() {
		result[cn] = true
	}
	cn.GetAllInterfacesWithSet(result)
	return result
}

func (cn *ClassNode) GetAllInterfacesWithSet(set map[IClassNode]bool) {
	for _, face := range cn.GetInterfaces() {
		if _, exists := set[face]; !exists {
			set[face] = true
			face.GetAllInterfacesWithSet(set)
		}
	}
}

func (cn *ClassNode) GetAbstractMethods() []MethodOrConstructorNode {
	r := cn.Redirect()
	var abstractMethods []MethodOrConstructorNode

	for _, methods := range r.methods {
		for _, method := range methods {
			if method.IsAbstract() {
				abstractMethods = append(abstractMethods, method)
			}
		}
	}

	return abstractMethods
}

func Values(m map[string]MethodOrConstructorNode) []MethodOrConstructorNode {
	values := make([]MethodOrConstructorNode, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetDeclaredMethodsFromSuper(cNode *ClassNode) map[string]MethodOrConstructorNode {
	parent := cNode.GetSuperClass()
	if parent == nil {
		return make(map[string]MethodOrConstructorNode)
	}
	return parent.GetDeclaredMethodsMap()
}

func (cn *ClassNode) GetAllDeclaredMethods() []MethodOrConstructorNode {
	return Values(cn.GetDeclaredMethodsMap())
}

// AddDeclaredMethodsFromInterfaces adds methods from all interfaces.
// Existing entries in the methods map take precedence. Methods from interfaces
// visited early take precedence over later ones.
//
// Parameters:
//   - cNode: The ClassNode
//   - methodsMap: A map of existing methods to alter
func AddDeclaredMethodsFromInterfaces(cNode *ClassNode, methodsMap map[string]MethodOrConstructorNode) {
	for _, iface := range cNode.GetInterfaces() {
		declaredMethods := iface.GetDeclaredMethodsMap()
		for key, method := range declaredMethods {
			if method.GetDeclaringClass().IsInterface() && (method.GetModifiers()&ACC_SYNTHETIC) == 0 {
				if _, exists := methodsMap[key]; !exists {
					methodsMap[key] = method
				}
			}
		}
	}
}

func (cn *ClassNode) GetDeclaredMethodsMap() map[string]MethodOrConstructorNode {
	result := GetDeclaredMethodsFromSuper(cn)
	AddDeclaredMethodsFromInterfaces(cn, result)

	// Add methods implemented in this class
	for _, methods := range cn.methods {
		for _, method := range methods {
			result[method.GetTypeDescriptor()] = method
		}
	}

	return result
}

func (cn *ClassNode) GetAnnotations() []*AnnotationNode {
	if cn.redirect != nil {
		return cn.redirect.GetAnnotations()
	}
	cn.lazyClassInit()
	return cn.AnnotatedNode.GetAnnotations()
}

func (cn *ClassNode) GetAnnotationsWithType(type_ IClassNode) []*AnnotationNode {
	if cn.redirect != nil {
		return cn.redirect.GetAnnotationsWithType(type_)
	}
	cn.lazyClassInit()
	return cn.AnnotatedNode.GetAnnotationsOfType(type_)
}
