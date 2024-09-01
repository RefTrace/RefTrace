package outputs

import "go.starlark.net/starlark"

type Output interface {
	starlark.Value
	starlark.HasAttrs
}
