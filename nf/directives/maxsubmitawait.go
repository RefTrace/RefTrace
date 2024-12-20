package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (m *MaxSubmitAwaitDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(m.Line()),
		Directive: &pb.Directive_MaxSubmitAwait{
			MaxSubmitAwait: &pb.MaxSubmitAwaitDirective{
				MaxSubmitAwait: m.MaxSubmitAwait,
			},
		},
	}
}

var _ Directive = (*MaxSubmitAwaitDirective)(nil)
var _ starlark.Value = (*MaxSubmitAwaitDirective)(nil)
var _ starlark.HasAttrs = (*MaxSubmitAwaitDirective)(nil)

func (m *MaxSubmitAwaitDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "max_submit_await":
		return starlark.String(m.MaxSubmitAwait), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("max submit await directive has no attribute %q", name))
	}
}

func (m *MaxSubmitAwaitDirective) AttrNames() []string {
	return []string{"max_submit_await"}
}

type MaxSubmitAwaitDirective struct {
	MaxSubmitAwait string
	line           int
}

func (m *MaxSubmitAwaitDirective) Line() int {
	return m.line
}

func (m *MaxSubmitAwaitDirective) String() string {
	return fmt.Sprintf("MaxSubmitAwaitDirective(MaxSubmitAwait: %q)", m.MaxSubmitAwait)
}

func (m *MaxSubmitAwaitDirective) Type() string {
	return "max_submit_await_directive"
}

func (m *MaxSubmitAwaitDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MaxSubmitAwaitDirective) Truth() starlark.Bool {
	return starlark.Bool(m.MaxSubmitAwait != "")
}

func (m *MaxSubmitAwaitDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(m.MaxSubmitAwait))
	return h.Sum32(), nil
}

func MakeMaxSubmitAwaitDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max submit await directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &MaxSubmitAwaitDirective{MaxSubmitAwait: strValue, line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid max submit await directive")
}
