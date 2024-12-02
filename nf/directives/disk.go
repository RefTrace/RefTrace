package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"

	pb "reft-go/nf/proto"
)

func (d *DiskDirective) ToProto() *pb.Directive {
	return &pb.Directive{
		Line: int32(d.Line()),
		Directive: &pb.Directive_Disk{
			Disk: &pb.DiskDirective{
				Space: d.Space,
			},
		},
	}
}

var _ Directive = (*DiskDirective)(nil)

func (d *DiskDirective) String() string {
	return fmt.Sprintf("DiskDirective(Space: %q)", d.Space)
}

func (d *DiskDirective) Type() string {
	return "disk_directive"
}

func (d *DiskDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (d *DiskDirective) Truth() starlark.Bool {
	return starlark.Bool(d.Space != "")
}

func (d *DiskDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(d.Space))
	return h.Sum32(), nil
}

var _ starlark.Value = (*DiskDirective)(nil)
var _ starlark.HasAttrs = (*DiskDirective)(nil)

func (d *DiskDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "space":
		return starlark.String(d.Space), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("disk directive has no attribute %q", name))
	}
}

func (d *DiskDirective) AttrNames() []string {
	return []string{"space"}
}

type DiskDirective struct {
	Space string
	line  int
}

func (d *DiskDirective) Line() int {
	return d.line
}

func MakeDiskDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid disk directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &DiskDirective{Space: strValue, line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid disk directive")
}
