package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*StoreDirDirective)(nil)

type StoreDirDirective struct {
	Directory string
}

func (a StoreDirDirective) Type() DirectiveType { return StoreDirDirectiveType }

func MakeStoreDirDirective(mce *parser.MethodCallExpression) (*StoreDirDirective, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid StoreDir directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &StoreDirDirective{Directory: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid StoreDir directive")
}
