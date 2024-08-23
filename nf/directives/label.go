package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*LabelDirective)(nil)

type LabelDirective struct {
	Label string
}

func (a LabelDirective) Type() DirectiveType { return LabelDirectiveType }

func MakeLabelDirective(mce *parser.MethodCallExpression) (*LabelDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid label directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &LabelDirective{Label: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid label directive")
}
