package parser

import "strings"

type ClassNode struct {
	BaseASTNode
	DefaultNodeMetaDataHandler
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
}

func NewClassNode(name string, modifiers int, superClass *ClassNode) *ClassNode {
	return &ClassNode{
		name:       name,
		modifiers:  modifiers,
		superClass: superClass,
		methods:    make(map[string][]*MethodNode),
	}
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
