package directives

import (
	"fmt"
	"hash/fnv"

	"go.starlark.net/starlark"
)

type DirectiveType int

const (
	AcceleratorType DirectiveType = iota
	AfterScriptType
	ArchType
	ArrayDirectiveType
	BeforeScriptType
	CacheDirectiveType
	ClusterOptionsType
	CondaType
	ContainerType
	ContainerOptionsType
	CpusDirectiveType
	DebugDirectiveType
	DiskDirectiveType
	EchoDirectiveType
	ErrorStrategyDirectiveType
	ExecutorDirectiveType
	ExtDirectiveType
	FairDirectiveType
	LabelDirectiveType
	MachineTypeDirectiveType
	MaxSubmitAwaitDirectiveType
	MaxErrorsDirectiveType
	MaxForksDirectiveType
	MaxRetriesDirectiveType
	MemoryDirectiveType
	ModuleDirectiveType
	PenvDirectiveType
	PodDirectiveType
	PublishDirDirectiveType
	QueueDirectiveType
	ResourceLabelsDirectiveType
	ResourceLimitsDirectiveType
	ScratchDirectiveType
	ShellDirectiveType
	SpackDirectiveType
	StageInModeDirectiveType
	StageOutModeDirectiveType
	StoreDirDirectiveType
	TagDirectiveType
	TimeDirectiveType
	DynamicDirectiveType
	UnknownDirectiveType
)

type Directive interface {
	starlark.Value
}

var _ Directive = (*DynamicDirective)(nil)
var _ starlark.Value = (*DynamicDirective)(nil)
var _ starlark.HasAttrs = (*DynamicDirective)(nil)

func (d *DynamicDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(d.Name), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("dynamic directive has no attribute %q", name))
	}
}

func (d *DynamicDirective) AttrNames() []string {
	return []string{"name"}
}

type DynamicDirective struct {
	Name string
}

func (d *DynamicDirective) String() string {
	return fmt.Sprintf("DynamicDirective(Name: %q)", d.Name)
}

func (d *DynamicDirective) Type() string {
	return "dynamic_directive"
}

func (d *DynamicDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (d *DynamicDirective) Truth() starlark.Bool {
	return starlark.Bool(d.Name != "")
}

func (d *DynamicDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(d.Name))
	return h.Sum32(), nil
}

var _ Directive = (*UnknownDirective)(nil)
var _ starlark.Value = (*UnknownDirective)(nil)
var _ starlark.HasAttrs = (*UnknownDirective)(nil)

func (u *UnknownDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(u.Name), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("unknown directive has no attribute %q", name))
	}
}

func (u *UnknownDirective) AttrNames() []string {
	return []string{"name"}
}

type UnknownDirective struct {
	Name string
}

func (u *UnknownDirective) String() string {
	return fmt.Sprintf("UnknownDirective(Name: %q)", u.Name)
}

func (u *UnknownDirective) Type() string {
	return "unknown_directive"
}

func (u *UnknownDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (u *UnknownDirective) Truth() starlark.Bool {
	return starlark.Bool(u.Name != "")
}

func (u *UnknownDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(u.Name))
	return h.Sum32(), nil
}
