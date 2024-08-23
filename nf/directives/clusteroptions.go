package directives

import (
	"reft-go/parser"
	"strings"
)

var _ Directive = (*ClusterOptions)(nil)

type ClusterOptions struct {
	Options string
}

func (a ClusterOptions) Type() DirectiveType { return ClusterOptionsType }

func MakeClusterOptions(mce *parser.MethodCallExpression) *ClusterOptions {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		options := []string{}
		for _, expr := range exprs {
			if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
				text := constantExpr.GetText()
				options = append(options, text)
			}
		}
		joinedOptions := strings.Join(options, " ")
		return &ClusterOptions{Options: joinedOptions}
	}
	return nil
}
