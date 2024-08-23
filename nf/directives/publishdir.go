package directives

import (
	"errors"
	"reft-go/parser"
)

var _ Directive = (*PublishDirDirective)(nil)

type PublishDirDirective struct {
	Path        string
	ContentType *bool
	Enabled     *bool
	FailOnError *bool
	Mode        string
	Overwrite   *bool
}

func (a PublishDirDirective) Type() DirectiveType { return PublishDirDirectiveType }

func MakePublishDirDirective(mce *parser.MethodCallExpression) (*PublishDirDirective, error) {
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
