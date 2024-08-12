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
	Classes           []IClassNode
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
	ScriptDummy       IClassNode
	ImportsResolved   bool
}

func NewModuleNode(description string) *ModuleNode {
	// TODO: initialize the base AST node
	return &ModuleNode{
		BaseASTNode:       NewBaseASTNode(),
		Description:       description,
		Classes:           []IClassNode{},
		Methods:           []*MethodNode{},
		Imports:           []*ImportNode{},
		StarImports:       []*ImportNode{},
		StaticImports:     make(map[string]*ImportNode),
		StaticStarImports: make(map[string]*ImportNode),
		StatementBlock:    NewBlockStatement(),
	}
}

func (m *ModuleNode) GetClasses() []IClassNode {
	mainClass := m.createStatementsClass()
	m.MainClassName = mainClass.GetName()
	m.Classes = append([]IClassNode{mainClass}, m.Classes...)
	//mainClass.SetModule(m)
	//m.AddToCompileUnit(mainClass)
	return m.Classes // Note: Go doesn't have built-in immutable collections
}

func (m *ModuleNode) createStatementsClass() IClassNode {
	classNode := m.GetScriptClassDummy()
	if strings.HasSuffix(classNode.GetName(), "package-info") {
		return classNode
	}

	//hasUncontainedStatements := false
	var fields []*FieldNode

	// Check for uncontained statements (excluding decl statements)
	for _, statement := range m.StatementBlock.GetStatements() {
		if _, ok := statement.(*ExpressionStatement); !ok {
			//hasUncontainedStatements = true
			break
		}
		es := statement.(*ExpressionStatement)
		expression := es.GetExpression()
		if _, ok := expression.(*DeclarationExpression); !ok {
			//hasUncontainedStatements = true
			break
		}
		de := expression.(*DeclarationExpression)
		if de.IsMultipleAssignmentDeclaration() {
			variables := de.GetTupleExpression().GetExpressions()
			rightExpr, ok := de.GetRightExpression().(*ListExpression)
			if !ok {
				break
			}
			values := rightExpr.GetExpressions()
			for i, v := range variables {
				varExpr := v.(*VariableExpression)
				var val Expression
				if i < len(values) {
					val = values[i]
				}
				fields = append(fields, NewFieldNode(varExpr.GetName(), varExpr.GetModifiers(), varExpr.GetType(), nil, val))
			}
		} else {
			ve := de.GetVariableExpression()
			fields = append(fields, NewFieldNode(ve.GetName(), ve.GetModifiers(), ve.GetType(), nil, de.GetRightExpression()))
		}
	}

	methodNode := NewMethodNode("run", ACC_PUBLIC, OBJECT_TYPE, []*Parameter{}, []IClassNode{}, m.StatementBlock)
	methodNode.SetIsScriptBody()
	AddGeneratedMethod(classNode, methodNode, true)

	classNode.AddConstructorWithDetails(ACC_PUBLIC, []*Parameter{}, []IClassNode{}, NewBlockStatement())

	var stmt Statement

	stmt, _ = NewExpressionStatement(NewConstructorCallExpression(
		SUPER_EXPRESSION.GetType(),
		NewArgumentListExpressionFromSlice(NewVariableExpressionWithString("context")),
	))

	classNode.AddConstructorWithDetails(ACC_PUBLIC, []*Parameter{NewParameter(BINDING_TYPE, "context")}, []IClassNode{}, stmt)

	for _, method := range m.Methods {
		if method.IsAbstract() {
			panic(fmt.Sprintf("Cannot use abstract methods in a script, they are only available inside classes. Method: %s", method.GetName()))
		}
		classNode.AddMethod(method)
	}

	return classNode
}

func AddGeneratedMethod(cNode IClassNode, mNode *MethodNode, skipChecks bool) {
	cNode.AddMethod(mNode)
	// TODO: implement MarkAsGenerated
	// AnnotatedNodeUtils.MarkAsGenerated(cNode, mNode, skipChecks)
}

func (m *ModuleNode) setScriptBaseClassFromConfig(cn IClassNode) {
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

			var superClass IClassNode
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

func (m *ModuleNode) GetScriptClassDummy() IClassNode {
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

	var classNode IClassNode

	classNode = NewClassNode(name, ACC_PUBLIC, SCRIPT_TYPE)
	m.setScriptBaseClassFromConfig(classNode)
	classNode.SetScript(true)
	classNode.SetScriptBody(true)

	m.ScriptDummy = classNode
	return classNode
}

func (mn *ModuleNode) Contains(class IClassNode) bool {
	for _, c := range mn.Classes {
		if c.Equals(class) {
			return true
		}
	}
	return false
}

func (m *ModuleNode) AddClass(node IClassNode) {
	if len(m.Classes) == 0 {
		m.MainClassName = node.GetName()
	}
	m.Classes = append(m.Classes, node)
	// TODO: check this
	//node.Module = m
	//m.AddToCompileUnit(node)
}

func (m *ModuleNode) AddToCompileUnit(node IClassNode) {
	if node != nil {
		m.Unit.AddClass(node)
	}
}

func (m *ModuleNode) AddMethod(node *MethodNode) {
	m.Methods = append(m.Methods, node)
}

func (m *ModuleNode) AddImportWithAnnotations(name string, classNode IClassNode, annotations []*AnnotationNode) {
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

func (m *ModuleNode) AddImport(name string, classNode IClassNode) {
	importNode := NewImportNodeType(classNode, name)
	m.Imports = append(m.Imports, importNode)
}

func (m *ModuleNode) AddStarImport(packageName string) {
	importNode := NewImportNodePackage(packageName)
	m.StarImports = append(m.StarImports, importNode)
}

func (m *ModuleNode) AddStaticImport(classNode IClassNode, memberName, simpleName string, annotations []*AnnotationNode) {
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

func (m *ModuleNode) AddStaticStarImport(name string, classNode IClassNode) {
	m.AddStaticStarImportWithAnnotations(name, classNode, nil)
}

func (m *ModuleNode) AddStaticStarImportWithAnnotations(name string, classNode IClassNode, annotations []*AnnotationNode) {
	importNode := NewImportNodeStatic(classNode)
	importNode.AddAnnotations(annotations)
	m.StaticStarImports[name] = importNode

	m.storeLastAddedImportNode(importNode)
}

func (m *ModuleNode) AddStatement(statement Statement) {
	m.StatementBlock.AddStatement(statement)
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

func (m *ModuleNode) checkUsage(name string, typ IClassNode) {
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
