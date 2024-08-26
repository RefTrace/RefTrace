package directives

import "go.starlark.net/starlark"

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

type DynamicDirective struct {
	Name string
}

func (a DynamicDirective) DirectiveType() DirectiveType { return DynamicDirectiveType }

var _ Directive = (*UnknownDirective)(nil)

type UnknownDirective struct {
	Name string
}

func (a UnknownDirective) DirectiveType() DirectiveType { return UnknownDirectiveType }
