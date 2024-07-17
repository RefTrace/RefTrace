package parser

import "strings"

var EMPTY_ARRAY = []*ClassNode{}
var THIS = OBJECT_TYPE
var SUPER = OBJECT_TYPE

type ClassNode struct {
	AnnotatedNode
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
}

func NewClassNode(name string, modifiers int, superClass *ClassNode) *ClassNode {
	return &ClassNode{
		name:       name,
		modifiers:  modifiers,
		superClass: superClass,
		methods:    make(map[string][]*MethodNode),
		fieldIndex: make(map[string]*FieldNode),
	}
}

func NewClassNodeWithInterfaces(name string, modifiers int, superClass *ClassNode, interfaces []*ClassNode, mixins []*MixinNode) *ClassNode {
	cn := NewClassNode(name, modifiers, superClass)
	cn.interfaces = interfaces
	// Assuming you want to add mixins to the ClassNode, you might need to add a mixins field to the ClassNode struct
	// cn.mixins = mixins
	return cn
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

func (cn *ClassNode) IsStatic() bool {
	return (cn.modifiers & ACC_STATIC) != 0
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
