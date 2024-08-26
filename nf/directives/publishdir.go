package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*PublishDirDirective)(nil)
var _ starlark.Value = (*PublishDirDirective)(nil)
var _ starlark.HasAttrs = (*PublishDirDirective)(nil)

func (p *PublishDirDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "path":
		return starlark.String(p.Path), nil
	case "contentType":
		if p.ContentType != nil {
			return starlark.Bool(*p.ContentType), nil
		}
		return starlark.None, nil
	case "enabled":
		if p.Enabled != nil {
			return starlark.Bool(*p.Enabled), nil
		}
		return starlark.None, nil
	case "failOnError":
		if p.FailOnError != nil {
			return starlark.Bool(*p.FailOnError), nil
		}
		return starlark.None, nil
	case "mode":
		return starlark.String(p.Mode), nil
	case "overwrite":
		if p.Overwrite != nil {
			return starlark.Bool(*p.Overwrite), nil
		}
		return starlark.None, nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("publish_dir directive has no attribute %q", name))
	}
}

func (p *PublishDirDirective) AttrNames() []string {
	return []string{"path", "contentType", "enabled", "failOnError", "mode", "overwrite"}
}

type PublishDirDirective struct {
	Path        string
	ContentType *bool
	Enabled     *bool
	FailOnError *bool
	Mode        string
	Overwrite   *bool
}

func (p *PublishDirDirective) String() string {
	return fmt.Sprintf("PublishDirDirective(Path: %q, ContentType: %v, Enabled: %v, FailOnError: %v, Mode: %q, Overwrite: %v)",
		p.Path, boolPtrToString(p.ContentType), boolPtrToString(p.Enabled), boolPtrToString(p.FailOnError), p.Mode, boolPtrToString(p.Overwrite))
}

func (p *PublishDirDirective) Type() string {
	return "publish_dir_directive"
}

func (p *PublishDirDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (p *PublishDirDirective) Truth() starlark.Bool {
	return starlark.Bool(p.Path != "")
}

func (p *PublishDirDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(p.Path))
	h.Write([]byte(boolPtrToString(p.ContentType)))
	h.Write([]byte(boolPtrToString(p.Enabled)))
	h.Write([]byte(boolPtrToString(p.FailOnError)))
	h.Write([]byte(p.Mode))
	h.Write([]byte(boolPtrToString(p.Overwrite)))
	return h.Sum32(), nil
}

// Helper function to convert *bool to string
func boolPtrToString(b *bool) string {
	if b == nil {
		return "nil"
	}
	return fmt.Sprintf("%t", *b)
}

func MakePublishDirDirective(mce *parser.MethodCallExpression) (Directive, error) {
	var dir string = ""
	var contentType *bool = nil
	var enabled *bool = nil
	var failOnError *bool = nil
	var mode string = ""
	var overwrite *bool = nil
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		for _, expr := range exprs {
			if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if strValue, ok := value.(string); ok {
					dir = strValue
				}
			}
			if mapExpr, ok := expr.(*parser.MapExpression); ok {
				entries := mapExpr.GetMapEntryExpressions()
				for _, entry := range entries {
					if entry.GetKeyExpression().GetText() == "contentType" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							v := constantExpr.GetValue()
							if vb, ok := v.(bool); ok {
								contentType = &vb
							}
						}
					}
					if entry.GetKeyExpression().GetText() == "enabled" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							v := constantExpr.GetValue()
							if vb, ok := v.(bool); ok {
								enabled = &vb
							}
						}
					}
					if entry.GetKeyExpression().GetText() == "failOnError" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							v := constantExpr.GetValue()
							if vb, ok := v.(bool); ok {
								failOnError = &vb
							}
						}
					}
					if entry.GetKeyExpression().GetText() == "mode" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							mode = constantExpr.GetText()
						}
					}
					if entry.GetKeyExpression().GetText() == "overwrite" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							v := constantExpr.GetValue()
							if vb, ok := v.(bool); ok {
								overwrite = &vb
							}
						}
					}
					if entry.GetKeyExpression().GetText() == "path" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							dir = constantExpr.GetText()
						}
					}
				}
			}
		}
	}
	if dir != "" {
		return &PublishDirDirective{
			Path:        dir,
			ContentType: contentType,
			Enabled:     enabled,
			FailOnError: failOnError,
			Mode:        mode,
			Overwrite:   overwrite,
		}, nil
	}
	return nil, errors.New("invalid publish dir directive")
}
