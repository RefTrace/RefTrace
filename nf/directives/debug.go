package directives

import (
	"reft-go/parser"
)

var _ Directive = (*DebugDirective)(nil)

type DebugDirective struct {
	Enabled bool
}

func (a DebugDirective) Type() DirectiveType { return DebugDirectiveType }

func MakeDebugDirective(mce *parser.MethodCallExpression) *DebugDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if boolValue, ok := constantExpr.GetValue().(bool); ok {
				return &DebugDirective{Enabled: boolValue}
			}
		}
	}
	return nil
}
