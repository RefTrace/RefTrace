package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ScratchDirective)(nil)

type ScratchDirective struct {
	Enabled   bool
	Directory string
}

func (a ScratchDirective) Type() DirectiveType { return ScratchDirectiveType }

func MakeScratchDirective(mce *parser.MethodCallExpression) (Directive, error) {
	enabled := false
	directory := ""
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Scratch directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				enabled = true
				directory = strValue
			}
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				enabled = boolValue
			}
		}
		if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
			enabled = true
			directory = gstringExpr.GetText()
		}
	}
	return &ScratchDirective{Enabled: enabled, Directory: directory}, nil
}
