package parser

import (
	"hash/fnv"
)

// RecordComponentNode represents a record component
type RecordComponentNode struct {
	AnnotatedNode
	name           string
	componentType  *ClassNode
	declaringClass *ClassNode
}

// NewRecordComponentNode creates a new RecordComponentNode with the given declaringClass, name, and type
func NewRecordComponentNode(declaringClass *ClassNode, name string, componentType *ClassNode) *RecordComponentNode {
	return NewRecordComponentNodeWithAnnotations(declaringClass, name, componentType, nil)
}

// NewRecordComponentNodeWithAnnotations creates a new RecordComponentNode with the given declaringClass, name, type, and annotations
func NewRecordComponentNodeWithAnnotations(declaringClass *ClassNode, name string, componentType *ClassNode, annotations []*AnnotationNode) *RecordComponentNode {
	rcn := &RecordComponentNode{
		name:           name,
		componentType:  componentType,
		declaringClass: declaringClass,
	}
	for _, annotation := range annotations {
		rcn.AddAnnotationNode(annotation)
	}
	return rcn
}

// GetName returns the name of the record component
func (rcn *RecordComponentNode) GetName() string {
	return rcn.name
}

// GetType returns the type of the record component
func (rcn *RecordComponentNode) GetType() *ClassNode {
	return rcn.componentType
}

// GetDeclaringClass returns the declaring class of the record component
func (rcn *RecordComponentNode) GetDeclaringClass() *ClassNode {
	return rcn.declaringClass
}

// Equals checks if the RecordComponentNode is equal to another object
func (rcn *RecordComponentNode) Equals(o interface{}) bool {
	if rcn == o {
		return true
	}
	other, ok := o.(*RecordComponentNode)
	if !ok {
		return false
	}
	return rcn.name == other.name && rcn.declaringClass.Equals(other.declaringClass)
}

// HashCode returns the hash code of the RecordComponentNode
func (rcn *RecordComponentNode) HashCode() uint32 {
	h := fnv.New32a()
	h.Write([]byte(rcn.declaringClass.GetName()))
	h.Write([]byte(rcn.name))
	return h.Sum32()
}
