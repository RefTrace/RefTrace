package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*StageOutModeDirective)(nil)

type StageOutModeDirective struct {
	Mode string
}

func (a StageOutModeDirective) Type() DirectiveType { return StageOutModeDirectiveType }

func MakeStageOutModeDirective(mce *parser.MethodCallExpression) (*StageOutModeDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid StageOutMode directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &StageOutModeDirective{Mode: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid StageOutMode directive")
}
