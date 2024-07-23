package parser

import (
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
)

type ModuleNode struct {
	*BaseASTNode
	Classes           []*ClassNode
	Methods           []*MethodNode
	Imports           []*ImportNode
	StarImports       []*ImportNode
	StaticImports     map[string]*ImportNode
	StaticStarImports map[string]*ImportNode
	Unit              *CompileUnit
	PackageNode       *PackageNode
	Description       string
	MainClassName     string
	StatementBlock    *BlockStatement
	ScriptDummy       *ClassNode
	ImportsResolved   bool
}

func NewModuleNode() *ModuleNode {
	// TODO: initialize the base AST node
	return &ModuleNode{
		Classes:           []*ClassNode{},
		Methods:           []*MethodNode{},
		Imports:           []*ImportNode{},
		StarImports:       []*ImportNode{},
		StaticImports:     make(map[string]*ImportNode),
		StaticStarImports: make(map[string]*ImportNode),
	}
}

func (m *ModuleNode) setScriptBaseClassFromConfig(cn *ClassNode) {
	/*
		var baseClassName string = ""
		//var bcLoader ClassLoader

		if m.Unit != nil {
			//bcLoader = m.Unit.GetClassLoader()
			baseClassName = m.Unit.GetConfig().GetScriptBaseClass()
		} else if m.Context != nil {
			//bcLoader = m.Context.GetClassLoader()
			baseClassName = m.Context.GetConfiguration().GetScriptBaseClass()
		}


		if baseClassName != "" && cn.GetSuperClass().GetName() != baseClassName {
			cn.AddAnnotation(NewAnnotationNode(BaseScriptASTTransformation.MY_TYPE))

			var superClass *ClassNode
			loadedClass, err := bcLoader.LoadClass(baseClassName)
			if err == nil {
				superClass = Make(loadedClass)
			} else {
				// If loading fails, fall back to making the class without loading
				superClass = MakeFromString(baseClassName)
			}
			cn.SetSuperClass(superClass)
		}
	*/
	return
}

func (m *ModuleNode) extractClassFromFileDescription() string {
	answer := m.Description

	// Parse the URI
	uri, err := url.Parse(answer)
	if err == nil {
		path := uri.Path
		schemeSpecific := uri.Opaque // Go's equivalent to getSchemeSpecificPart
		if path != "" {
			answer = path
		} else if schemeSpecific != "" {
			answer = schemeSpecific
		}
	}

	// Strip off everything after the last '.'
	slashIdx := strings.LastIndex(answer, "/")
	separatorIdx := strings.LastIndex(answer, string(filepath.Separator))
	dotIdx := strings.LastIndex(answer, ".")
	if dotIdx > 0 && dotIdx > max(slashIdx, separatorIdx) {
		answer = answer[:dotIdx]
	}

	// Strip everything up to and including the path separators
	if slashIdx >= 0 {
		answer = answer[slashIdx+1:]
	}

	// Recalculate in case we have already done some stripping
	separatorIdx = strings.LastIndex(answer, string(filepath.Separator))
	if separatorIdx >= 0 {
		answer = answer[separatorIdx+1:]
	}

	return answer
}

// Helper function to find the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *ModuleNode) GetScriptClassDummy() *ClassNode {
	if m.ScriptDummy != nil {
		m.setScriptBaseClassFromConfig(m.ScriptDummy)
		return m.ScriptDummy
	}

	name := m.GetPackageName()
	if name == "" {
		name = ""
	}
	// now let's use the file name to determine the class name
	if m.Description == "" {
		panic("Cannot generate main(String[]) class for statements when we have no file description")
	}
	name += EncodeAsValidClassName(m.extractClassFromFileDescription())

	var classNode *ClassNode

	classNode = NewClassNode(name, ACC_PUBLIC, SCRIPT_TYPE)
	m.setScriptBaseClassFromConfig(classNode)
	classNode.SetScript(true)
	classNode.SetScriptBody(true)

	m.ScriptDummy = classNode
	return classNode
}

func (mn *ModuleNode) Contains(class *ClassNode) bool {
	for _, c := range mn.Classes {
		if c.Equals(class) {
			return true
		}
	}
	return false
}

func (m *ModuleNode) AddClass(node *ClassNode) {
	if len(m.Classes) == 0 {
		m.MainClassName = node.name
	}
	m.Classes = append(m.Classes, node)
	// TODO: check this
	//node.Module = m
	m.AddToCompileUnit(node)
}

func (m *ModuleNode) AddToCompileUnit(node *ClassNode) {
	if node != nil {
		m.Unit.AddClass(node)
	}
}

func (m *ModuleNode) AddMethod(node *MethodNode) {
	m.Methods = append(m.Methods, node)
}

func (m *ModuleNode) AddImportWithAnnotations(name string, classNode *ClassNode, annotations []*AnnotationNode) {
	importNode := NewImportNodeType(classNode, name)
	importNode.AddAnnotations(annotations)
	m.Imports = append(m.Imports, importNode)
	m.storeLastAddedImportNode(importNode)
}

func (m *ModuleNode) AddStarImportWithAnnotations(packageName string, annotations []*AnnotationNode) {
	importNode := NewImportNodePackage(packageName)
	importNode.AddAnnotations(annotations)
	m.StarImports = append(m.StarImports, importNode)
	m.storeLastAddedImportNode(importNode)
}

func (m *ModuleNode) AddImport(name string, classNode *ClassNode) {
	importNode := NewImportNodeType(classNode, name)
	m.Imports = append(m.Imports, importNode)
}

func (m *ModuleNode) AddStarImport(packageName string) {
	importNode := NewImportNodePackage(packageName)
	m.StarImports = append(m.StarImports, importNode)
}

func (m *ModuleNode) AddStaticImport(classNode *ClassNode, memberName, simpleName string, annotations []*AnnotationNode) {
	// Create a new ClassNode for the member
	memberType := NewClassNode(classNode.GetName()+"."+memberName, 0, nil)
	memberType.SetRedirect(classNode)
	memberType.SetSourcePosition(classNode)

	// Check usage
	m.checkUsage(simpleName, memberType)

	// Create the ImportNode
	importNode := NewImportNodeStaticField(classNode, memberName, simpleName)
	importNode.AddAnnotations(annotations)

	// Add to staticImports
	prev, exists := m.StaticImports[simpleName]
	if exists {
		m.StaticImports[prev.GetText()] = prev
		delete(m.StaticImports, simpleName)
	}
	m.StaticImports[simpleName] = importNode

	// Store last added import node
	m.storeLastAddedImportNode(importNode)
}

func (m *ModuleNode) storeLastAddedImportNode(node *ImportNode) {
	if m.GetNodeMetaData(reflect.TypeOf(&ImportNode{})) == reflect.TypeOf(&ImportNode{}) {
		m.PutNodeMetaData(reflect.TypeOf(&ImportNode{}), node)
	}
}

func (m *ModuleNode) AddStaticStarImport(name string, classNode *ClassNode) {
	m.AddStaticStarImportWithAnnotations(name, classNode, nil)
}

func (m *ModuleNode) AddStaticStarImportWithAnnotations(name string, classNode *ClassNode, annotations []*AnnotationNode) {
	importNode := NewImportNodeStatic(classNode)
	importNode.AddAnnotations(annotations)
	m.StaticStarImports[name] = importNode

	m.storeLastAddedImportNode(importNode)
}

func (m *ModuleNode) AddStatement(statement Statement) {
	m.StatementBlock.AddStatement(statement)
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

func (m *ModuleNode) GetStarImports() []*ImportNode {
	return m.StarImports
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

func (m *ModuleNode) GetStatementBlock() *BlockStatement {
	return m.StatementBlock
}

func (m *ModuleNode) GetPackageName() string {
	if m.PackageNode == nil {
		return ""
	}
	return m.PackageNode.name
}

func (m *ModuleNode) SetPackageName(packageName string) {
	m.PackageNode = NewPackageNode(packageName)
}

func (m *ModuleNode) checkUsage(name string, typ *ClassNode) {
	// Check classes
	for _, node := range m.Classes {
		if node.GetNameWithoutPackage() == name && !node.Equals(typ) {
			m.addErrorAndContinue(fmt.Errorf("the name %s is already declared", name), typ)
			return
		}
	}

	// Check imports
	for _, node := range m.Imports {
		if node.Alias == name && !node.Type.Equals(typ) {
			m.addErrorAndContinue(fmt.Errorf("the name %s is already declared", name), typ)
			return
		}
	}

	// Check static imports
	// TODO: properly implement this
	/*
		if node, exists := m.StaticImports[name]; exists {
			if !node.Type.Equals(typ.GetOuterClass()) {
				m.addErrorAndContinue(fmt.Errorf("the name %s is already declared", name), &typ.BaseASTNode)
			}
		}
	*/
}

// Helper function to add an error and continue
func (m *ModuleNode) addErrorAndContinue(err error, node ASTNode) {
	// You might want to implement a more sophisticated error handling mechanism
	// For now, we'll just print the error
	fmt.Printf("Error at line %v: %v\n", node.GetLineNumber(), err)
}
