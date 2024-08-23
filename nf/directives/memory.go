package directives

import (
	"reft-go/parser"
)

var _ Directive = (*MemoryDirective)(nil)

type MemoryDirective struct {
	Memory string
}

func (a MemoryDirective) Type() DirectiveType { return MemoryDirectiveType }

func MakeMemoryDirective(mce *parser.MethodCallExpression) *MemoryDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MemoryDirective{Memory: strValue}
			}
		}
	}
	return nil
}
