package directives

import (
	"errors"
	"reft-go/parser"
	"strings"
)

var _ Directive = (*Shell)(nil)

type Shell struct {
	Command string
}

func (a Shell) Type() DirectiveType { return ShellDirectiveType }

func MakeShellDirective(mce *parser.MethodCallExpression) (*Shell, error) {
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
		return &Shell{Command: joinedOptions}, nil
	}
	return nil, errors.New("invalid Shell directive")
}
