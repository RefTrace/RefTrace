package parser

import "reflect"

// var _ ASTNode = (*BaseASTNode)(nil)

type ASTNode interface {
	ASTNodeNoVisit
	Visit(visitor GroovyCodeVisitor)
}

type ASTNodeNoVisit interface {
	NodeMetaDataHandler
	GetText() string
	GetLineNumber() int
	SetLineNumber(int)
	GetColumnNumber() int
	SetColumnNumber(int)
	GetLastLineNumber() int
	SetLastLineNumber(int)
	GetLastColumnNumber() int
	SetLastColumnNumber(int)
	SetSourcePosition(ASTNode)
	GetMetaDataMap() map[interface{}]interface{}
	SetMetaDataMap(map[interface{}]interface{})
	CopyNodeMetaDataHandler(NodeMetaDataHandler)
	NewMetaDataMap() map[interface{}]interface{}
}

type BaseASTNode struct {
	*DefaultNodeMetaDataHandler
	lineNumber       int
	columnNumber     int
	lastLineNumber   int
	lastColumnNumber int
}

func (node *BaseASTNode) GetText() string {
	nodeType := reflect.TypeOf(node)
	if nodeType.Kind() == reflect.Ptr {
		nodeType = nodeType.Elem()
	}
	return "<not implemented yet for class: " + nodeType.String() + ">"
}

func (node *BaseASTNode) GetLineNumber() int {
	return node.lineNumber
}

func (node *BaseASTNode) SetLineNumber(lineNumber int) {
	node.lineNumber = lineNumber
}

func (node *BaseASTNode) GetColumnNumber() int {
	return node.columnNumber
}

func (node *BaseASTNode) SetColumnNumber(columnNumber int) {
	node.columnNumber = columnNumber
}

func (node *BaseASTNode) GetLastLineNumber() int {
	return node.lastLineNumber
}

func (node *BaseASTNode) SetLastLineNumber(lastLineNumber int) {
	node.lastLineNumber = lastLineNumber
}

func (node *BaseASTNode) GetLastColumnNumber() int {
	return node.lastColumnNumber
}

func (node *BaseASTNode) SetLastColumnNumber(lastColumnNumber int) {
	node.lastColumnNumber = lastColumnNumber
}

func (node *BaseASTNode) SetSourcePosition(other ASTNode) {
	node.lineNumber = other.GetLineNumber()
	node.columnNumber = other.GetColumnNumber()
	node.lastLineNumber = other.GetLastLineNumber()
	node.lastColumnNumber = other.GetLastColumnNumber()
}

func (node *BaseASTNode) GetMetaDataMap() map[interface{}]interface{} {
	return node.metaDataMap
}

func (node *BaseASTNode) SetMetaDataMap(metaDataMap map[interface{}]interface{}) {
	node.metaDataMap = metaDataMap
}

func (node *BaseASTNode) CopyNodeMetaDataHandler(handler NodeMetaDataHandler) {
	node.metaDataMap = handler.GetMetaDataMap()
}

func (node *BaseASTNode) NewMetaDataMap() map[interface{}]interface{} {
	node.metaDataMap = make(map[interface{}]interface{})
	return node.metaDataMap
}

// NewBaseASTNode creates and initializes a new BaseASTNode
func NewBaseASTNode() *BaseASTNode {
	return &BaseASTNode{
		DefaultNodeMetaDataHandler: &DefaultNodeMetaDataHandler{
			metaDataMap: make(map[interface{}]interface{}),
		},
		lineNumber:       0,
		columnNumber:     0,
		lastLineNumber:   0,
		lastColumnNumber: 0,
	}
}
