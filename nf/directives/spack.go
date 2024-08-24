package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*SpackDirective)(nil)

type SpackDirective struct {
	Dependencies string
}

func (a SpackDirective) Type() DirectiveType { return SpackDirectiveType }

func MakeSpackDirective(mce *parser.MethodCallExpression) (*SpackDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid Spack directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &SpackDirective{Dependencies: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid Spack directive")
}
