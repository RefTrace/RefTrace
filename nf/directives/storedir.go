package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*StoreDirDirective)(nil)
var _ starlark.Value = (*StoreDirDirective)(nil)
var _ starlark.HasAttrs = (*StoreDirDirective)(nil)

func (s *StoreDirDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "directory":
		return starlark.String(s.Directory), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("store_dir directive has no attribute %q", name))
	}
}

func (s *StoreDirDirective) AttrNames() []string {
	return []string{"directory"}
}

type StoreDirDirective struct {
	Directory string
	line      int
}

func (s *StoreDirDirective) Line() int {
	return s.line
}

func (s *StoreDirDirective) String() string {
	return fmt.Sprintf("StoreDirDirective(Directory: %q)", s.Directory)
}

func (s *StoreDirDirective) Type() string {
	return "store_dir_directive"
}

func (s *StoreDirDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (s *StoreDirDirective) Truth() starlark.Bool {
	return starlark.Bool(s.Directory != "")
}

func (s *StoreDirDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(s.Directory))
	return h.Sum32(), nil
}

func MakeStoreDirDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid StoreDir directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if strValue, ok := constantExpr.GetValue().(string); ok {
				return &StoreDirDirective{Directory: strValue, line: mce.GetLineNumber()}, nil
			}
		}
	}
	return nil, errors.New("invalid StoreDir directive")
}
