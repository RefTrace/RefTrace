package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*MaxSubmitAwaitDirective)(nil)

type MaxSubmitAwaitDirective struct {
	MaxSubmitAwait string
}

func (a MaxSubmitAwaitDirective) Type() DirectiveType { return MaxSubmitAwaitDirectiveType }

func MakeMaxSubmitAwaitDirective(mce *parser.MethodCallExpression) (*MaxSubmitAwaitDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max submit await directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MaxSubmitAwaitDirective{MaxSubmitAwait: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max submit await directive")
}
