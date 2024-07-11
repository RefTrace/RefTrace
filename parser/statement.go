package parser

import (
	"container/list"
)

// Statement represents the interface for any statement.
type Statement interface {
	ASTNode
	GetStatementLabels() *list.List
	GetStatementLabel() string
	GetText() string
	SetStatementLabel(label string)
	AddStatementLabel(label string)
	CopyStatementLabels(that Statement)
	ClearStatementLabels()
	IsEmpty() bool
}

// BaseStatement represents the base struct for any statement implementation.
type BaseStatement struct {
	ASTNode
	statementLabels *list.List
}

// GetStatementLabels returns the list of statement labels.
func (s *BaseStatement) GetStatementLabels() *list.List {
	return s.statementLabels
}

// GetStatementLabel returns the first statement label (deprecated).
func (s *BaseStatement) GetStatementLabel() string {
	if s.statementLabels == nil || s.statementLabels.Len() == 0 {
		return ""
	}
	return s.statementLabels.Front().Value.(string)
}

// SetStatementLabel sets a single statement label (deprecated).
func (s *BaseStatement) SetStatementLabel(label string) {
	if label != "" {
		s.AddStatementLabel(label)
	}
}

// AddStatementLabel adds a statement label to the list.
func (s *BaseStatement) AddStatementLabel(label string) {
	if s.statementLabels == nil {
		s.statementLabels = list.New()
	}
	s.statementLabels.PushBack(label)
}

// CopyStatementLabels copies statement labels from another Statement.
func (s *BaseStatement) CopyStatementLabels(that Statement) {
	if thatLabels := that.GetStatementLabels(); thatLabels != nil {
		for e := thatLabels.Front(); e != nil; e = e.Next() {
			s.AddStatementLabel(e.Value.(string))
		}
	}
}

// ClearStatementLabels removes all statement labels from the list.
func (s *BaseStatement) ClearStatementLabels() {
	if s.statementLabels != nil {
		s.statementLabels.Init()
	}
}

// IsEmpty returns whether the statement is empty.
func (s *BaseStatement) IsEmpty() bool {
	return false
}
