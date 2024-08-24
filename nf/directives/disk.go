package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*DiskDirective)(nil)

type DiskDirective struct {
	Space string
}

func (a DiskDirective) Type() DirectiveType { return DiskDirectiveType }

func MakeDiskDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid disk directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &DiskDirective{Space: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid disk directive")
}
