package parser

import (
	"fmt"
)

// Embed BaseASTNode in all node types
type ClassNode struct {
	BaseASTNode
	Name       string
	SuperClass *ClassNode
	Methods    []*MethodNode
	Fields     []*FieldNode
}

type MethodNode struct {
	BaseASTNode
	Name       string
	IsStatic   bool
	IsAbstract bool
	Parameters []string
	ReturnType string
}

type FieldNode struct {
	BaseASTNode
	Name string
	Type string
}

type ImportNode struct {
	BaseASTNode
	Name string
	Type *ClassNode
}

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

type Statement struct {
	BaseASTNode
	Content string
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

func main() {
	module := NewModuleNode()
	classNode := &ClassNode{Name: "ExampleClass"}
	module.AddClass(classNode)
	module.AddMethod(&MethodNode{Name: "exampleMethod", IsStatic: true, ReturnType: "void"})
	module.AddImport("exampleImport", classNode)
	module.AddStaticImport("exampleStaticImport", classNode)
	module.AddStaticStarImport("exampleStaticStarImport", classNode)
	module.AddStatement(&Statement{Content: "example statement"})

	fmt.Println("Main Class Name:", module.GetMainClassName())
	fmt.Println("Classes:", module.GetClasses())
	fmt.Println("Methods:", module.GetMethods())
	fmt.Println("Imports:", module.GetImports())
	fmt.Println("Static Imports:", module.GetStaticImports())
	fmt.Println("Static Star Imports:", module.GetStaticStarImports())
	fmt.Println("Statement Block:", module.GetStatementBlock())
}
