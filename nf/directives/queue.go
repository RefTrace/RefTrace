package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*QueueDirective)(nil)

type QueueDirective struct {
	Name string
}

func (a QueueDirective) Type() DirectiveType { return QueueDirectiveType }

func MakeQueueDirective(mce *parser.MethodCallExpression) (*QueueDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Queue directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &QueueDirective{Name: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid Queue directive")
}
