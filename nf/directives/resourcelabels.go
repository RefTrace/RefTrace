package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*ResourceLabelsDirective)(nil)

type ResourceLabelsDirective struct {
	Keys []string
}

func (a ResourceLabelsDirective) Type() DirectiveType { return ResourceLabelsDirectiveType }

func MakeResourceLabelsDirective(mce *parser.MethodCallExpression) (*ResourceLabelsDirective, error) {
	var keys []string = []string{}
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		if len(args.GetExpressions()) != 1 {
			return nil, errors.New("invalid resource labels directive")
		}
		expr := args.GetExpressions()[0]
		if namedArgListExpr, ok := expr.(*parser.NamedArgumentListExpression); ok {
			exprs := namedArgListExpr.GetMapEntryExpressions()
			for _, mapEntryExpr := range exprs {
				key := mapEntryExpr.GetKeyExpression().GetText()
				keys = append(keys, key)
			}
		}
	}
	return &ResourceLabelsDirective{Keys: keys}, nil
}
