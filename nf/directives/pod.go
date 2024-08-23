package directives

import (
	"reft-go/parser"
)

var _ Directive = (*PodDirective)(nil)

type PodDirective struct {
	Env   string
	Value string
}

func (a PodDirective) Type() DirectiveType { return PodDirectiveType }

func MakePodDirective(mce *parser.MethodCallExpression) *PodDirective {
	var env string = ""
	var val string = ""
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		if len(args.GetExpressions()) != 1 {
			return nil
		}
		expr := args.GetExpressions()[0]
		if namedArgListExpr, ok := expr.(*parser.NamedArgumentListExpression); ok {
			exprs := namedArgListExpr.GetMapEntryExpressions()
			for _, mapEntryExpr := range exprs {

				key := mapEntryExpr.GetKeyExpression().GetText()
				value := mapEntryExpr.GetValueExpression()
				if key == "env" {
					if constantExpr, ok := value.(*parser.ConstantExpression); ok {
						env = constantExpr.GetText()
					}
				}
				if key == "value" {
					if constantExpr, ok := value.(*parser.ConstantExpression); ok {
						val = constantExpr.GetText()
					}
				}

			}
		}
	}
	if env != "" && val != "" {
		return &PodDirective{Env: env, Value: val}
	}
	return nil
}
