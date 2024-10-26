package parser

// AnnotatedNode represents an AST node that can be annotated
type AnnotatedNode struct {
	*BaseASTNode
	annotations    []*AnnotationNode
	declaringClass IClassNode
	synthetic      bool
}

// NewAnnotatedNode creates and initializes a new AnnotatedNode
func NewAnnotatedNode() *AnnotatedNode {
	return &AnnotatedNode{
		BaseASTNode:    NewBaseASTNode(),
		annotations:    make([]*AnnotationNode, 0),
		declaringClass: nil,
		synthetic:      false,
	}
}

// GetAnnotations returns all annotations for this node
func (an *AnnotatedNode) GetAnnotations() []*AnnotationNode {
	return an.annotations
}

// GetAnnotationsOfType returns annotations of a specific type
func (an *AnnotatedNode) GetAnnotationsOfType(typ IClassNode) []*AnnotationNode {
	var result []*AnnotationNode
	for _, node := range an.annotations {
		if node.GetClassNode().Equals(typ) {
			result = append(result, node)
		}
	}
	return result
}

// AddAnnotation adds a new annotation of the given type
func (an *AnnotatedNode) AddAnnotation(typ IClassNode) *AnnotationNode {
	node := NewAnnotationNode(typ)
	an.AddAnnotationNode(node)
	return node
}

// AddAnnotationNode addAnnotationNode adds an existing annotation node
func (an *AnnotatedNode) AddAnnotationNode(annotation *AnnotationNode) {
	if annotation != nil {
		an.annotations = append(an.annotations, annotation)
	}
}

// AddAnnotations adds multiple annotations
func (an *AnnotatedNode) AddAnnotations(annotations []*AnnotationNode) {
	for _, annotation := range annotations {
		an.AddAnnotationNode(annotation)
	}
}

// GetDeclaringClass returns the declaring class of this node
func (an *AnnotatedNode) GetDeclaringClass() IClassNode {
	return an.declaringClass
}

// SetDeclaringClass sets the declaring class of this node
func (an *AnnotatedNode) SetDeclaringClass(declaringClass IClassNode) {
	an.declaringClass = declaringClass
}

// HasNoRealSourcePosition returns true for default constructors added by the compiler
func (an *AnnotatedNode) HasNoRealSourcePosition() bool {
	val, ok := an.GetNodeMetaData("org.codehaus.groovy.ast.AnnotatedNode.hasNoRealSourcePosition").(bool)
	return ok && val
}

// SetHasNoRealSourcePosition sets whether this node has no real source position
func (an *AnnotatedNode) SetHasNoRealSourcePosition(hasNoRealSourcePosition bool) {
	if hasNoRealSourcePosition {
		an.PutNodeMetaData("org.codehaus.groovy.ast.AnnotatedNode.hasNoRealSourcePosition", true)
	} else {
		an.RemoveNodeMetaData("org.codehaus.groovy.ast.AnnotatedNode.hasNoRealSourcePosition")
	}
}

// IsSynthetic indicates if this node was added by the compiler
func (an *AnnotatedNode) IsSynthetic() bool {
	return an.synthetic
}

// SetSynthetic sets this node as a node added by the compiler
func (an *AnnotatedNode) SetSynthetic(synthetic bool) {
	an.synthetic = synthetic
}
