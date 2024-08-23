package directives

import (
	"reft-go/parser"
)

var _ Directive = (*DiskDirective)(nil)

type DiskDirective struct {
	Space string
}

func (a DiskDirective) Type() DirectiveType { return DiskDirectiveType }

func MakeDiskDirective(mce *parser.MethodCallExpression) *DiskDirective {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &DiskDirective{Space: strValue}
			}
		}
	}
	return nil
}
