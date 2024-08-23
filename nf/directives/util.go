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
)

type Directive interface {
	Type() DirectiveType
}
