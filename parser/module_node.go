package parser

type ModuleNode struct {
	BaseASTNode
	Classes           []*ClassNode
	Methods           []*MethodNode
	Imports           []*ImportNode
	StaticImports     map[string]*ImportNode
	StaticStarImports map[string]*ImportNode
	MainClassName     string
	StatementBlock    []*Statement
}

func NewModuleNode() *ModuleNode {
	return &ModuleNode{
		Classes:           []*ClassNode{},
		Methods:           []*MethodNode{},
		Imports:           []*ImportNode{},
		StaticImports:     make(map[string]*ImportNode),
		StaticStarImports: make(map[string]*ImportNode),
		StatementBlock:    []*Statement{},
	}
}

func (m *ModuleNode) AddClass(node *ClassNode) {
	if len(m.Classes) == 0 {
		m.MainClassName = node.Name
	}
	m.Classes = append(m.Classes, node)
}

func (m *ModuleNode) AddMethod(node *MethodNode) {
	m.Methods = append(m.Methods, node)
}

func (m *ModuleNode) AddImport(name string, classNode *ClassNode) {
	importNode := &ImportNode{Name: name, Type: classNode}
	m.Imports = append(m.Imports, importNode)
}

func (m *ModuleNode) AddStaticImport(name string, classNode *ClassNode) {
	importNode := &ImportNode{Name: name, Type: classNode}
	m.StaticImports[name] = importNode
}

func (m *ModuleNode) AddStaticStarImport(name string, classNode *ClassNode) {
	importNode := &ImportNode{Name: name, Type: classNode}
	m.StaticStarImports[name] = importNode
}

func (m *ModuleNode) AddStatement(statement *Statement) {
	m.StatementBlock = append(m.StatementBlock, statement)
}

func (m *ModuleNode) GetClasses() []*ClassNode {
	return m.Classes
}

func (m *ModuleNode) GetMethods() []*MethodNode {
	return m.Methods
}

func (m *ModuleNode) GetImports() []*ImportNode {
	return m.Imports
}

func (m *ModuleNode) GetStaticImports() map[string]*ImportNode {
	return m.StaticImports
}

func (m *ModuleNode) GetStaticStarImports() map[string]*ImportNode {
	return m.StaticStarImports
}

func (m *ModuleNode) GetMainClassName() string {
	return m.MainClassName
}

func (m *ModuleNode) GetStatementBlock() []*Statement {
	return m.StatementBlock
}
