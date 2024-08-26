package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*ExtDirective)(nil)

func (e *ExtDirective) String() string {
	return fmt.Sprintf("ExtDirective(Version: %q, Args: %q)", e.Version, e.Args)
}

func (e *ExtDirective) Type() string {
	return "ext_directive"
}

func (e *ExtDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (e *ExtDirective) Truth() starlark.Bool {
	return starlark.Bool(e.Version != "" || e.Args != "")
}

func (e *ExtDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(e.Version))
	h.Write([]byte(e.Args))
	return h.Sum32(), nil
}

type ExtDirective struct {
	Version string
	Args    string
}

func MakeExtDirective(mce *parser.MethodCallExpression) (Directive, error) {
	var version string = ""
	var extArgs string = ""
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		if len(args.GetExpressions()) != 1 {
			return nil, errors.New("invalid ext directive")
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
		return &ExtDirective{Version: version, Args: extArgs}, nil
	}
	return nil, errors.New("invalid ext directive")
}
