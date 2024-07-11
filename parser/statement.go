package parser

import (
	"container/list"
)

// Statement represents the base struct for any statement.
type Statement struct {
	ASTNode
	statementLabels *list.List
}

// GetStatementLabels returns the list of statement labels.
func (s *Statement) GetStatementLabels() *list.List {
	return s.statementLabels
}

// GetStatementLabel returns the first statement label (deprecated).
func (s *Statement) GetStatementLabel() string {
	if s.statementLabels == nil || s.statementLabels.Len() == 0 {
		return ""
	}
	return s.statementLabels.Front().Value.(string)
}

// SetStatementLabel sets a single statement label (deprecated).
func (s *Statement) SetStatementLabel(label string) {
	if label != "" {
		s.AddStatementLabel(label)
	}
}

// AddStatementLabel adds a statement label to the list.
func (s *Statement) AddStatementLabel(label string) {
	if s.statementLabels == nil {
		s.statementLabels = list.New()
	}
	s.statementLabels.PushBack(label)
}

// CopyStatementLabels copies statement labels from another Statement.
func (s *Statement) CopyStatementLabels(that *Statement) {
	if that.statementLabels != nil {
		for e := that.statementLabels.Front(); e != nil; e = e.Next() {
			s.AddStatementLabel(e.Value.(string))
		}
	}
}

// ClearStatementLabels removes all statement labels from the list.
func (s *Statement) ClearStatementLabels() {
	if s.statementLabels != nil {
		s.statementLabels.Init()
	}
}

// IsEmpty returns whether the statement is empty.
func (s *Statement) IsEmpty() bool {
	return false
}
