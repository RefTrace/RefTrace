package nf

import (
	"fmt"
	"reft-go/parser"
	"strings"
)

var _ parser.IClassNode = (*UnionTypeClassNode)(nil)

type UnionTypeClassNode struct {
	parser.IClassNode
	delegates []parser.IClassNode
}

func NewUnionTypeClassNode(classNodes ...parser.IClassNode) *UnionTypeClassNode {
	base := parser.NewClassNode(makeName(classNodes), 0, makeSuper(classNodes))
	return &UnionTypeClassNode{
		IClassNode: base,
		delegates:  classNodes,
	}
}

func makeName(nodes []parser.IClassNode) string {
	names := make([]string, len(nodes))
	for i, node := range nodes {
		names[i] = node.GetName()
	}
	return fmt.Sprintf("<UnionType:%s>", strings.Join(names, "+"))
}

func makeSuper(nodes []parser.IClassNode) parser.IClassNode {
	upper := parser.LowestUpperBound(nodes)
	if _, ok := upper.(*parser.LowestUpperBoundClassNode); ok {
		upper = upper.GetUnresolvedSuperClass()
	} else if upper.IsInterface() {
		upper = parser.OBJECT_TYPE
	}
	return upper
}

func (u *UnionTypeClassNode) GetDelegates() []parser.IClassNode {
	return u.delegates
}

// Implement other methods from ClassNode interface
// Most of these methods will either aggregate results from delegates
// or throw UnsupportedOperationError

func (u *UnionTypeClassNode) AddConstructor(node *parser.ConstructorNode) {
	panic("UnsupportedOperationError")
}

func (u *UnionTypeClassNode) AddField(node *parser.FieldNode) {
	panic("UnsupportedOperationError")
}

func (u *UnionTypeClassNode) AddInterface(type_ parser.IClassNode) {
	panic("UnsupportedOperationError")
}

func (u *UnionTypeClassNode) AddMethod(node parser.MethodOrConstructorNode) {
	panic("UnsupportedOperationError")
}

func (u *UnionTypeClassNode) GetAbstractMethods() []parser.MethodOrConstructorNode {
	var methods []parser.MethodOrConstructorNode
	for _, delegate := range u.delegates {
		methods = append(methods, delegate.GetAbstractMethods()...)
	}
	return methods
}

func (u *UnionTypeClassNode) GetAllDeclaredMethods() []parser.MethodOrConstructorNode {
	var methods []parser.MethodOrConstructorNode
	for _, delegate := range u.delegates {
		methods = append(methods, delegate.GetAllDeclaredMethods()...)
	}
	return methods
}

func (u *UnionTypeClassNode) GetAllInterfaces() map[parser.IClassNode]bool {
	interfaces := make(map[parser.IClassNode]bool)
	for _, delegate := range u.delegates {
		for iface, value := range delegate.GetAllInterfaces() {
			interfaces[iface] = value
		}
	}
	return interfaces
}

func (u *UnionTypeClassNode) GetAnnotations() []*parser.AnnotationNode {
	var annotations []*parser.AnnotationNode
	for _, delegate := range u.delegates {
		annotations = append(annotations, delegate.GetAnnotations()...)
	}
	return annotations
}

func (u *UnionTypeClassNode) GetDeclaredField(name string) *parser.FieldNode {
	for _, delegate := range u.delegates {
		if field := delegate.GetDeclaredField(name); field != nil {
			return field
		}
	}
	return nil
}

func (u *UnionTypeClassNode) GetInterfaces() []parser.IClassNode {
	interfaces := make(map[parser.IClassNode]bool)
	for _, delegate := range u.delegates {
		if delegate.IsInterface() {
			interfaces[delegate] = true
		} else {
			for _, iface := range delegate.GetInterfaces() {
				interfaces[iface] = true
			}
		}
	}
	result := make([]parser.IClassNode, 0, len(interfaces))
	for iface := range interfaces {
		result = append(result, iface)
	}
	return result
}

func (u *UnionTypeClassNode) IsDerivedFrom(type_ parser.IClassNode) bool {
	for _, delegate := range u.delegates {
		if delegate.IsDerivedFrom(type_) {
			return true
		}
	}
	return false
}
