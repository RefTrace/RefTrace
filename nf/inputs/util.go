package inputs

import "go.starlark.net/starlark"

type Input interface {
	starlark.Value
	starlark.HasAttrs
}
