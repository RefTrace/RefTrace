package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

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

type DiskDirective struct {
	Space string
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
				return &DiskDirective{Space: strValue}, nil
			}
		}
	}
	return nil, errors.New("invalid disk directive")
}
