package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*PodDirective)(nil)

type PodDirective struct {
	Env   string
	Value string
}

func (p *PodDirective) String() string {
	return fmt.Sprintf("PodDirective(Env: %q, Value: %q)", p.Env, p.Value)
}

func (p *PodDirective) Type() string {
	return "pod_directive"
}

func (p *PodDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (p *PodDirective) Truth() starlark.Bool {
	return starlark.Bool(p.Env != "" && p.Value != "")
}

func (p *PodDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(p.Env))
	h.Write([]byte(p.Value))
	return h.Sum32(), nil
}

func MakePodDirective(mce *parser.MethodCallExpression) (Directive, error) {
	var env string = ""
	var val string = ""
	if args, ok := mce.GetArguments().(*parser.TupleExpression); ok {
		if len(args.GetExpressions()) != 1 {
			return nil, errors.New("invalid pod directive")
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
		return &PodDirective{Env: env, Value: val}, nil
	}
	return nil, errors.New("invalid pod directive")
}
