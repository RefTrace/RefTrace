package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*StageInModeDirective)(nil)

type StageInModeDirective struct {
	Mode string
}

func (a StageInModeDirective) Type() DirectiveType { return StageInModeDirectiveType }

func MakeStageInModeDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid StageInMode directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &StageInModeDirective{Mode: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid StageInMode directive")
}
