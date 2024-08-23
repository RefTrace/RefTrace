package directives

import (
	"reft-go/parser"
)

var _ Directive = (*ExtDirective)(nil)

type ExtDirective struct {
	Version string
	Args    string
}

func (a ExtDirective) Type() DirectiveType { return ExtDirectiveType }

func MakeExtDirective(mce *parser.MethodCallExpression) *ExtDirective {
	var version string = ""
	var extArgs string = ""
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
				if key == "version" {
					if constantExpr, ok := value.(*parser.ConstantExpression); ok {
						version = constantExpr.GetText()
					}
				}
				if key == "args" {
					if constantExpr, ok := value.(*parser.ConstantExpression); ok {
						extArgs = constantExpr.GetText()
					}
				}

			}
		}
	}
	if version != "" {
		return &ExtDirective{Version: version, Args: extArgs}
	}
	return nil
}
