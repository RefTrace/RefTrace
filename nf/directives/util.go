package directives

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
	Type() DirectiveType
}

var _ Directive = (*DynamicDirective)(nil)

type DynamicDirective struct {
	Name string
}

func (a DynamicDirective) Type() DirectiveType { return DynamicDirectiveType }

var _ Directive = (*UnknownDirective)(nil)

type UnknownDirective struct {
	Name string
}

func (a UnknownDirective) Type() DirectiveType { return UnknownDirectiveType }
