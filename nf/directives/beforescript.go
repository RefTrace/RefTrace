package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (b *BeforeScript) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(b.Line()),
		Directive: &pb.Directive_BeforeScript{
			BeforeScript: &pb.BeforeScriptDirective{
				Script: b.Script,
			},
		},
	}
}

var _ Directive = (*BeforeScript)(nil)
var _ starlark.Value = (*BeforeScript)(nil)
var _ starlark.HasAttrs = (*BeforeScript)(nil)

func (b *BeforeScript) Attr(name string) (starlark.Value, error) {
	switch name {
	case "script":
		return starlark.String(b.Script), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("beforescript directive has no attribute %q", name))
	}
}

func (b *BeforeScript) AttrNames() []string {
	return []string{"script"}
}

func (b *BeforeScript) String() string {
	return fmt.Sprintf("BeforeScript(%q)", b.Script)
}

func (b *BeforeScript) Type() string {
	return "before_script"
}

func (b *BeforeScript) Freeze() {
	// No mutable fields, so no action needed
}

func (b *BeforeScript) Truth() starlark.Bool {
	return starlark.Bool(b.Script != "")
}

func (b *BeforeScript) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(b.Script))
	return h.Sum32(), nil
}

type BeforeScript struct {
	Script string
	line   int
}

func (b *BeforeScript) Line() int {
	return b.line
}

func MakeBeforeScript(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) == 1 {
			if constantExpr, ok := exprs[0].(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					return &BeforeScript{Script: strValue, line: mce.GetLineNumber()}, nil
				}
			}
			if gstringExpr, ok := exprs[0].(*parser.GStringExpression); ok {
				return &BeforeScript{Script: gstringExpr.GetText(), line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid beforeScript directive")
}
