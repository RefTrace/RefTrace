package parser

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

var _ antlr.ErrorStrategy = &CustomErrorStrategy{}

type CustomErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	errors []error
}

// CustomError wraps a RecognitionException to implement the error interface
type CustomError struct {
	antlr.RecognitionException
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("Recognition error: %s", ce.GetMessage())
}

func NewCustomErrorStrategy() *CustomErrorStrategy {
	return &CustomErrorStrategy{
		DefaultErrorStrategy: antlr.NewDefaultErrorStrategy(),
		errors:               []error{},
	}
}

func (c *CustomErrorStrategy) ReportError(_ antlr.Parser, e antlr.RecognitionException) {
	c.errors = append(c.errors, &CustomError{e})
}

func (c *CustomErrorStrategy) ReportMatch(recognizer antlr.Parser) {
	c.DefaultErrorStrategy.ReportMatch(recognizer)
}

func (c *CustomErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	c.errors = append(c.errors, &CustomError{e})
	c.DefaultErrorStrategy.Recover(recognizer, e)
}

func (c *CustomErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	c.errors = append(c.errors, &CustomError{antlr.NewInputMisMatchException(recognizer)})
	return c.DefaultErrorStrategy.RecoverInline(recognizer)
}

func (c *CustomErrorStrategy) Sync(recognizer antlr.Parser) {
	c.DefaultErrorStrategy.Sync(recognizer)
}

func (c *CustomErrorStrategy) GetErrors() []error {
	return c.errors
}

func (c *CustomErrorStrategy) ClearErrors() {
	c.errors = []error{}
}

func (c *CustomErrorStrategy) ReportNoViableAlternative(recognizer antlr.Parser, e *antlr.NoViableAltException) {
	c.errors = append(c.errors, &CustomError{e})
}

func (c *CustomErrorStrategy) ReportInputMisMatch(recognizer antlr.Parser, e *antlr.InputMisMatchException) {
	c.errors = append(c.errors, &CustomError{e})
}

func (c *CustomErrorStrategy) ReportFailedPredicate(recognizer antlr.Parser, e *antlr.FailedPredicateException) {
	c.errors = append(c.errors, &CustomError{e})
}

func (c *CustomErrorStrategy) ReportUnwantedToken(recognizer antlr.Parser) {
	e := antlr.NewInputMisMatchException(recognizer)
	c.errors = append(c.errors, &CustomError{e})
}

func (c *CustomErrorStrategy) ReportMissingToken(recognizer antlr.Parser) {
	e := antlr.NewInputMisMatchException(recognizer)
	c.errors = append(c.errors, &CustomError{e})
}
