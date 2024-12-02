package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (e *ExtDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(e.Line()),
		Directive: &pb.Directive_Ext{
			Ext: &pb.ExtDirective{
				Version: e.Version,
				Args:    e.Args,
			},
		},
	}
}

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

var _ starlark.Value = (*ExtDirective)(nil)
var _ starlark.HasAttrs = (*ExtDirective)(nil)

func (e *ExtDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "version":
		return starlark.String(e.Version), nil
	case "args":
		return starlark.String(e.Args), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("ext directive has no attribute %q", name))
	}
}

func (e *ExtDirective) AttrNames() []string {
	return []string{"version", "args"}
}

type ExtDirective struct {
	Version string
	Args    string
	line    int
}

func (e *ExtDirective) Line() int {
	return e.line
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
		return &ExtDirective{Version: version, Args: extArgs, line: mce.GetLineNumber()}, nil
	}
	return nil, errors.New("invalid ext directive")
}
