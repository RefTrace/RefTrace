package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (p *PodDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(p.Line()),
		Directive: &pb.Directive_Pod{
			Pod: &pb.PodDirective{
				Env:   p.Env,
				Value: p.Value,
			},
		},
	}
}

var _ Directive = (*PodDirective)(nil)
var _ starlark.Value = (*PodDirective)(nil)
var _ starlark.HasAttrs = (*PodDirective)(nil)

func (p *PodDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "env":
		return starlark.String(p.Env), nil
	case "value":
		return starlark.String(p.Value), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("pod directive has no attribute %q", name))
	}
}

func (p *PodDirective) AttrNames() []string {
	return []string{"env", "value"}
}

type PodDirective struct {
	Env   string
	Value string
	line  int
}

func (p *PodDirective) Line() int {
	return p.line
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
		return &PodDirective{Env: env, Value: val, line: mce.GetLineNumber()}, nil
	}
	return nil, errors.New("invalid pod directive")
}
