package parser

import (
	"sync"
)

// CompileUnit represents the entire contents of a compilation step
type CompileUnit struct {
	metaDataMap           map[interface{}]interface{}
	modules               []*ModuleNode
	classes               map[string]*ClassNode
	classesToCompile      map[string]*ClassNode
	generatedInnerClasses map[string]*InnerClassNode
	mu                    sync.RWMutex
}

// NewCompileUnit creates a new CompileUnit
func NewCompileUnit() *CompileUnit {
	return &CompileUnit{
		classes:               make(map[string]*ClassNode),
		classesToCompile:      make(map[string]*ClassNode),
		generatedInnerClasses: make(map[string]*InnerClassNode),
	}
}

// GetMetaDataMap returns the metadata map
func (cu *CompileUnit) GetMetaDataMap() map[interface{}]interface{} {
	return cu.metaDataMap
}

// SetMetaDataMap sets the metadata map
func (cu *CompileUnit) SetMetaDataMap(metaDataMap map[interface{}]interface{}) {
	cu.metaDataMap = metaDataMap
}

// GetModules returns the list of ModuleNodes
func (cu *CompileUnit) GetModules() []*ModuleNode {
	return cu.modules
}

// GetClasses returns a list of all classes in each module
func (cu *CompileUnit) GetClasses() []*ClassNode {
	cu.mu.RLock()
	defer cu.mu.RUnlock()

	var answer []*ClassNode
	for _, module := range cu.modules {
		answer = append(answer, module.GetClasses()...)
	}
	return answer
}

// GetClass returns the ClassNode for the given qualified name
func (cu *CompileUnit) GetClass(name string) *ClassNode {
	cu.mu.RLock()
	defer cu.mu.RUnlock()

	cn, ok := cu.classes[name]
	if !ok {
		cn, _ = cu.classesToCompile[name]
	}
	return cn
}

// GetClassesToCompile returns the map of classes to compile
func (cu *CompileUnit) GetClassesToCompile() map[string]*ClassNode {
	cu.mu.RLock()
	defer cu.mu.RUnlock()
	return cu.classesToCompile
}

// GetGeneratedInnerClasses returns the map of generated inner classes
func (cu *CompileUnit) GetGeneratedInnerClasses() map[string]*InnerClassNode {
	cu.mu.RLock()
	defer cu.mu.RUnlock()
	return cu.generatedInnerClasses
}

// GetGeneratedInnerClass returns the InnerClassNode for the given name
func (cu *CompileUnit) GetGeneratedInnerClass(name string) *InnerClassNode {
	cu.mu.RLock()
	defer cu.mu.RUnlock()
	return cu.generatedInnerClasses[name]
}

// AddClasses adds a list of ClassNodes to the CompileUnit
func (cu *CompileUnit) AddClasses(list []*ClassNode) {
	for _, node := range list {
		cu.AddClass(node)
	}
}

// AddClass adds a ClassNode to the CompileUnit
func (cu *CompileUnit) AddClass(node *ClassNode) {
	cu.mu.Lock()
	defer cu.mu.Unlock()

	node = node.Redirect()
	name := node.GetName()
	stored, ok := cu.classes[name]
	if ok && stored != node {
		// Handle duplicate class definition error
		// (Error handling logic would go here)
	}
	cu.classes[name] = node

	cn, ok := cu.classesToCompile[name]
	if ok {
		delete(cu.classesToCompile, name)
		cn.SetRedirect(node)
	}
}

// AddGeneratedInnerClass adds an InnerClassNode to the CompileUnit
func (cu *CompileUnit) AddGeneratedInnerClass(icn *InnerClassNode) {
	cu.mu.Lock()
	defer cu.mu.Unlock()
	cu.generatedInnerClasses[icn.GetName()] = icn
}
