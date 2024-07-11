package parser

import (
	"sync"
)

// VariableScope records declared and referenced variables for a given scope.
// It helps determine variable sharing across closure and method boundaries.
type VariableScope struct {
	parent                   *VariableScope
	classScope               *ClassNode
	inStaticContext          bool
	declaredVariables        map[string]*Variable
	referencedLocalVariables map[string]*Variable
	referencedClassVariables map[string]*Variable
	mu                       sync.RWMutex
}

func NewVariableScope() *VariableScope {
	return &VariableScope{
		declaredVariables:        make(map[string]*Variable),
		referencedLocalVariables: make(map[string]*Variable),
		referencedClassVariables: make(map[string]*Variable),
	}
}

func NewVariableScopeWithParent(parent *VariableScope) *VariableScope {
	vs := NewVariableScope()
	vs.parent = parent
	return vs
}

func (vs *VariableScope) GetParent() *VariableScope {
	return vs.parent
}

func (vs *VariableScope) IsRoot() bool {
	return vs.parent == nil
}

func (vs *VariableScope) GetClassScope() *ClassNode {
	return vs.classScope
}

func (vs *VariableScope) IsClassScope() bool {
	return vs.classScope != nil
}

func (vs *VariableScope) SetClassScope(classScope *ClassNode) {
	vs.classScope = classScope
}

func (vs *VariableScope) IsInStaticContext() bool {
	return vs.inStaticContext
}

func (vs *VariableScope) SetInStaticContext(inStaticContext bool) {
	vs.inStaticContext = inStaticContext
}

func (vs *VariableScope) GetDeclaredVariable(name string) *Variable {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.declaredVariables[name]
}

func (vs *VariableScope) GetReferencedLocalVariable(name string) *Variable {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.referencedLocalVariables[name]
}

func (vs *VariableScope) GetReferencedClassVariable(name string) *Variable {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return vs.referencedClassVariables[name]
}

func (vs *VariableScope) IsReferencedLocalVariable(name string) bool {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	_, exists := vs.referencedLocalVariables[name]
	return exists
}

func (vs *VariableScope) IsReferencedClassVariable(name string) bool {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	_, exists := vs.referencedClassVariables[name]
	return exists
}

func (vs *VariableScope) GetDeclaredVariables() map[string]*Variable {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	result := make(map[string]*Variable)
	for k, v := range vs.declaredVariables {
		result[k] = v
	}
	return result
}

func (vs *VariableScope) GetReferencedClassVariables() map[string]*Variable {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	result := make(map[string]*Variable)
	for k, v := range vs.referencedClassVariables {
		result[k] = v
	}
	return result
}

func (vs *VariableScope) GetReferencedLocalVariablesCount() int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return len(vs.referencedLocalVariables)
}

func (vs *VariableScope) PutDeclaredVariable(v *Variable) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.declaredVariables[(*v).Name()] = v
}

func (vs *VariableScope) PutReferencedLocalVariable(v *Variable) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.referencedLocalVariables[(*v).Name()] = v
}

func (vs *VariableScope) PutReferencedClassVariable(v *Variable) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.referencedClassVariables[(*v).Name()] = v
}

func (vs *VariableScope) RemoveReferencedClassVariable(name string) *Variable {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	v, exists := vs.referencedClassVariables[name]
	if exists {
		delete(vs.referencedClassVariables, name)
		return v
	}
	return nil
}

func (vs *VariableScope) Copy() *VariableScope {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	that := NewVariableScopeWithParent(vs.parent)
	that.classScope = vs.classScope
	that.inStaticContext = vs.inStaticContext

	for k, v := range vs.declaredVariables {
		that.declaredVariables[k] = v
	}
	for k, v := range vs.referencedLocalVariables {
		that.referencedLocalVariables[k] = v
	}
	for k, v := range vs.referencedClassVariables {
		that.referencedClassVariables[k] = v
	}

	return that
}
