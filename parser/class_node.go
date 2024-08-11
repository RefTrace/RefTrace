package parser

import (
	"fmt"
	"strings"
)

var EMPTY_ARRAY = []*ClassNode{}
var THIS = OBJECT_TYPE
var SUPER = OBJECT_TYPE

type ClassNode struct {
	*AnnotatedNode
	name                string
	modifiers           int
	superClass          *ClassNode
	interfaces          []*ClassNode
	fields              []*FieldNode
	methods             map[string][]*MethodNode
	constructors        []*ConstructorNode
	isPrimaryNode       bool
	redirect            *ClassNode
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

func NewClassNode(name string, modifiers int, superClass *ClassNode) *ClassNode {
	return &ClassNode{
		AnnotatedNode: NewAnnotatedNode(),
		name:          name,
		modifiers:     modifiers,
		superClass:    superClass,
		methods:       make(map[string][]*MethodNode),
		fieldIndex:    make(map[string]*FieldNode),
	}
}

func NewClassNodeWithInterfaces(name string, modifiers int, superClass *ClassNode, interfaces []*ClassNode, mixins []*MixinNode) *ClassNode {
	cn := NewClassNode(name, modifiers, superClass)
	cn.interfaces = interfaces
	// Assuming you want to add mixins to the ClassNode, you might need to add a mixins field to the ClassNode struct
	// cn.mixins = mixins
	return cn
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

func (cn *ClassNode) SetInterfaces(interfaces []*ClassNode) {
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

func (cn *ClassNode) GetOrAddStaticInitializer() *MethodNode {
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
		r.constructors = make([]*ConstructorNode, 0, 4)
	}
	r.constructors = append(r.constructors, node)
}

func (cn *ClassNode) AddConstructorWithDetails(modifiers int, parameters []*Parameter, exceptions []*ClassNode, code Statement) *ConstructorNode {
	node := NewConstructorNodeWithParams(modifiers, parameters, exceptions, code)
	cn.AddConstructor(node)
	return node
}

func (cn *ClassNode) GetDeclaredMethod(name string, parameters []*Parameter) *MethodNode {
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

func (cn *ClassNode) GetDeclaredMethods(name string) []*MethodNode {
	if cn.redirect != nil {
		return cn.redirect.GetDeclaredMethods(name)
	}
	return cn.methods[name]
}

func (cn *ClassNode) AddMethod(node *MethodNode) {
	r := cn.Redirect()
	node.SetDeclaringClass(r)
	if r.methods == nil {
		r.methods = make(map[string][]*MethodNode)
	}
	r.methods[node.GetName()] = append(r.methods[node.GetName()], node)
}

func (cn *ClassNode) AddMethodWithDetails(name string, modifiers int, returnType *ClassNode, parameters []*Parameter, exceptions []*ClassNode, code Statement) *MethodNode {
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

func (cn *ClassNode) IsDerivedFrom(type_ *ClassNode) bool {
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
	for node := cn; node != nil; node = node.GetSuperClass() {
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

func (cn *ClassNode) Equals(that interface{}) bool {
	other, ok := that.(*ClassNode)
	if !ok {
		return false
	}
	if cn == other {
		return true
	}
	if cn.redirect != nil {
		return cn.redirect.Equals(that)
	}
	if cn.componentType != nil {
		return cn.componentType.Equals(other.componentType)
	}
	return other.GetText() == cn.GetText() // arrays could be "T[]" or "[LT;"
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

func (cn *ClassNode) GetSuperClass() *ClassNode {
	if cn.redirect != nil {
		return cn.redirect.GetSuperClass()
	}
	return cn.superClass
}

func (cn *ClassNode) SetSuperClass(superClass *ClassNode) {
	if cn.redirect != nil {
		cn.redirect.SetSuperClass(superClass)
	} else {
		cn.superClass = superClass
	}
}

func (cn *ClassNode) GetInterfaces() []*ClassNode {
	if cn.redirect != nil {
		return cn.redirect.GetInterfaces()
	}
	return cn.interfaces
}

func (cn *ClassNode) AddInterface(iface *ClassNode) {
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
	n.SetRedirect(cn.redirect)

	if cn.IsArray() {
		n.componentType = cn.redirect.GetComponentType()
	}

	return n
}

func (cn *ClassNode) SetRedirect(node *ClassNode) {
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

func (cn *ClassNode) GetPlainNodeReference() *ClassNode {
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

func (cn *ClassNode) ImplementsInterface(classNode *ClassNode) bool {
	node := cn.Redirect()
	for node != nil {
		if node.DeclaresInterface(classNode) {
			return true
		}
		node = node.GetSuperClass()
	}
	return false
}

func (cn *ClassNode) DeclaresAnyInterfaces(classNodes ...*ClassNode) bool {
	for _, classNode := range classNodes {
		if cn.DeclaresInterface(classNode) {
			return true
		}
	}
	return false
}

func (cn *ClassNode) DeclaresInterface(classNode *ClassNode) bool {
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
		upper := cn
		if cn.redirect != nil {
			upper = cn.redirect
		}
		return NewGenericsType(cn, []*ClassNode{upper}, nil)
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

func (cn *ClassNode) GetOuterClass() *ClassNode {
	if cn.redirect != nil {
		return cn.redirect.GetOuterClass()
	}
	return nil
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
