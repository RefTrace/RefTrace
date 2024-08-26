package nf

import (
	"fmt"
	"reflect"
	"reft-go/nf/directives"

	"go.starlark.net/starlark"
)

func ConvertToStarlarkProcess(p Process) *StarlarkProcess {
	sp := &StarlarkProcess{
		Name:       p.Name,
		Directives: &StarlarkProcessDirectives{},
	}

	for _, directive := range p.Directives {
		switch d := directive.(type) {
		case *directives.Accelerator:
			sp.Directives.Accelerator = append(sp.Directives.Accelerator, d)
		case *directives.AfterScript:
			sp.Directives.AfterScript = append(sp.Directives.AfterScript, d)
		case *directives.Arch:
			sp.Directives.Arch = append(sp.Directives.Arch, d)
		case *directives.ArrayDirective:
			sp.Directives.Array = append(sp.Directives.Array, d)
		case *directives.BeforeScript:
			sp.Directives.BeforeScript = append(sp.Directives.BeforeScript, d)
		case *directives.CacheDirective:
			sp.Directives.Cache = append(sp.Directives.Cache, d)
		case *directives.ClusterOptions:
			sp.Directives.ClusterOptions = append(sp.Directives.ClusterOptions, d)
		case *directives.Conda:
			sp.Directives.Conda = append(sp.Directives.Conda, d)
		case *directives.Container:
			sp.Directives.Container = append(sp.Directives.Container, d)
		case *directives.ContainerOptions:
			sp.Directives.ContainerOptions = append(sp.Directives.ContainerOptions, d)
		case *directives.CpusDirective:
			sp.Directives.Cpus = append(sp.Directives.Cpus, d)
		case *directives.DebugDirective:
			sp.Directives.Debug = append(sp.Directives.Debug, d)
		case *directives.DiskDirective:
			sp.Directives.Disk = append(sp.Directives.Disk, d)
		case *directives.EchoDirective:
			sp.Directives.Echo = append(sp.Directives.Echo, d)
		case *directives.ErrorStrategyDirective:
			sp.Directives.ErrorStrategy = append(sp.Directives.ErrorStrategy, d)
		case *directives.ExecutorDirective:
			sp.Directives.Executor = append(sp.Directives.Executor, d)
		case *directives.ExtDirective:
			sp.Directives.Ext = append(sp.Directives.Ext, d)
		case *directives.FairDirective:
			sp.Directives.Fair = append(sp.Directives.Fair, d)
		case *directives.LabelDirective:
			sp.Directives.Label = append(sp.Directives.Label, d)
		case *directives.MachineTypeDirective:
			sp.Directives.MachineType = append(sp.Directives.MachineType, d)
		case *directives.MaxSubmitAwaitDirective:
			sp.Directives.MaxSubmitAwait = append(sp.Directives.MaxSubmitAwait, d)
		case *directives.MaxErrorsDirective:
			sp.Directives.MaxErrors = append(sp.Directives.MaxErrors, d)
		case *directives.MaxForksDirective:
			sp.Directives.MaxForks = append(sp.Directives.MaxForks, d)
		case *directives.MaxRetriesDirective:
			sp.Directives.MaxRetries = append(sp.Directives.MaxRetries, d)
		case *directives.MemoryDirective:
			sp.Directives.Memory = append(sp.Directives.Memory, d)
		case *directives.ModuleDirective:
			sp.Directives.Module = append(sp.Directives.Module, d)
		case *directives.PenvDirective:
			sp.Directives.Penv = append(sp.Directives.Penv, d)
		case *directives.PodDirective:
			sp.Directives.Pod = append(sp.Directives.Pod, d)
		case *directives.PublishDirDirective:
			sp.Directives.PublishDir = append(sp.Directives.PublishDir, d)
		case *directives.QueueDirective:
			sp.Directives.Queue = append(sp.Directives.Queue, d)
		case *directives.ResourceLabelsDirective:
			sp.Directives.ResourceLabels = append(sp.Directives.ResourceLabels, d)
		case *directives.ResourceLimitsDirective:
			sp.Directives.ResourceLimits = append(sp.Directives.ResourceLimits, d)
		case *directives.ScratchDirective:
			sp.Directives.Scratch = append(sp.Directives.Scratch, d)
		case *directives.Shell:
			sp.Directives.Shell = append(sp.Directives.Shell, d)
		case *directives.SpackDirective:
			sp.Directives.Spack = append(sp.Directives.Spack, d)
		case *directives.StageInModeDirective:
			sp.Directives.StageInMode = append(sp.Directives.StageInMode, d)
		case *directives.StageOutModeDirective:
			sp.Directives.StageOutMode = append(sp.Directives.StageOutMode, d)
		case *directives.StoreDirDirective:
			sp.Directives.StoreDir = append(sp.Directives.StoreDir, d)
		case *directives.TagDirective:
			sp.Directives.Tag = append(sp.Directives.Tag, d)
		case *directives.TimeDirective:
			sp.Directives.Time = append(sp.Directives.Time, d)
		case *directives.DynamicDirective:
			sp.Directives.Dynamic = append(sp.Directives.Dynamic, d)
		case *directives.UnknownDirective:
			sp.Directives.Unknown = append(sp.Directives.Unknown, d)
		}
	}

	return sp
}

var _ starlark.Value = (*StarlarkProcess)(nil)
var _ starlark.HasAttrs = (*StarlarkProcess)(nil)

type StarlarkProcess struct {
	Name       string
	Directives *StarlarkProcessDirectives
}

func (p *StarlarkProcess) AttrNames() []string {
	return []string{"name", "directives"}
}

type StarlarkProcessDirectives struct {
	Accelerator      []*directives.Accelerator
	AfterScript      []*directives.AfterScript
	Arch             []*directives.Arch
	Array            []*directives.ArrayDirective
	BeforeScript     []*directives.BeforeScript
	Cache            []*directives.CacheDirective
	ClusterOptions   []*directives.ClusterOptions
	Conda            []*directives.Conda
	Container        []*directives.Container
	ContainerOptions []*directives.ContainerOptions
	Cpus             []*directives.CpusDirective
	Debug            []*directives.DebugDirective
	Disk             []*directives.DiskDirective
	Echo             []*directives.EchoDirective
	ErrorStrategy    []*directives.ErrorStrategyDirective
	Executor         []*directives.ExecutorDirective
	Ext              []*directives.ExtDirective
	Fair             []*directives.FairDirective
	Label            []*directives.LabelDirective
	MachineType      []*directives.MachineTypeDirective
	MaxSubmitAwait   []*directives.MaxSubmitAwaitDirective
	MaxErrors        []*directives.MaxErrorsDirective
	MaxForks         []*directives.MaxForksDirective
	MaxRetries       []*directives.MaxRetriesDirective
	Memory           []*directives.MemoryDirective
	Module           []*directives.ModuleDirective
	Penv             []*directives.PenvDirective
	Pod              []*directives.PodDirective
	PublishDir       []*directives.PublishDirDirective
	Queue            []*directives.QueueDirective
	ResourceLabels   []*directives.ResourceLabelsDirective
	ResourceLimits   []*directives.ResourceLimitsDirective
	Scratch          []*directives.ScratchDirective
	Shell            []*directives.Shell
	Spack            []*directives.SpackDirective
	StageInMode      []*directives.StageInModeDirective
	StageOutMode     []*directives.StageOutModeDirective
	StoreDir         []*directives.StoreDirDirective
	Tag              []*directives.TagDirective
	Time             []*directives.TimeDirective
	Dynamic          []*directives.DynamicDirective
	Unknown          []*directives.UnknownDirective
}

func (p *StarlarkProcess) String() string {
	return fmt.Sprintf("Process(%s)", p.Name)
}

func (p *StarlarkProcess) Type() string {
	return "process"
}

func (p *StarlarkProcess) Freeze() {
	// No mutable fields, so no action needed
}

func (p *StarlarkProcess) Truth() starlark.Bool {
	return starlark.Bool(p.Name != "")
}

func (p *StarlarkProcess) Hash() (uint32, error) {
	// Implement hash function if needed
	return 0, fmt.Errorf("unhashable type: process")
}

func (p *StarlarkProcess) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(p.Name), nil
	case "directives":
		return &StarlarkProcessDirectivesWrapper{p.Directives}, nil
	default:
		return nil, fmt.Errorf("process has no attribute %q", name)
	}
}

var _ starlark.Value = (*StarlarkProcessDirectivesWrapper)(nil)
var _ starlark.HasAttrs = (*StarlarkProcessDirectivesWrapper)(nil)

type StarlarkProcessDirectivesWrapper struct {
	*StarlarkProcessDirectives
}

func (w *StarlarkProcessDirectivesWrapper) String() string {
	return "ProcessDirectives"
}

func (w *StarlarkProcessDirectivesWrapper) Type() string {
	return "process_directives"
}

func (w *StarlarkProcessDirectivesWrapper) Freeze() {
	// No mutable fields, so no action needed
}

func (w *StarlarkProcessDirectivesWrapper) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (w *StarlarkProcessDirectivesWrapper) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: process_directives")
}

func starlarkListFromDirectives(directives interface{}) *starlark.List {
	v := reflect.ValueOf(directives)
	if v.Kind() != reflect.Slice {
		return starlark.NewList(nil)
	}

	elements := make([]starlark.Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		elements[i] = v.Index(i).Interface().(starlark.Value)
	}

	return starlark.NewList(elements)
}

func (w *StarlarkProcessDirectivesWrapper) Attr(name string) (starlark.Value, error) {
	switch name {
	case "accelerator":
		return starlarkListFromDirectives(w.Accelerator), nil
	case "after_script":
		return starlarkListFromDirectives(w.AfterScript), nil
	case "arch":
		return starlarkListFromDirectives(w.Arch), nil
	case "array":
		return starlarkListFromDirectives(w.Array), nil
	case "before_script":
		return starlarkListFromDirectives(w.BeforeScript), nil
	case "cache":
		return starlarkListFromDirectives(w.Cache), nil
	case "cluster_options":
		return starlarkListFromDirectives(w.ClusterOptions), nil
	case "conda":
		return starlarkListFromDirectives(w.Conda), nil
	case "container":
		return starlarkListFromDirectives(w.Container), nil
	case "container_options":
		return starlarkListFromDirectives(w.ContainerOptions), nil
	case "cpus":
		return starlarkListFromDirectives(w.Cpus), nil
	case "debug":
		return starlarkListFromDirectives(w.Debug), nil
	case "disk":
		return starlarkListFromDirectives(w.Disk), nil
	case "echo":
		return starlarkListFromDirectives(w.Echo), nil
	case "error_strategy":
		return starlarkListFromDirectives(w.ErrorStrategy), nil
	case "executor":
		return starlarkListFromDirectives(w.Executor), nil
	case "ext":
		return starlarkListFromDirectives(w.Ext), nil
	case "fair":
		return starlarkListFromDirectives(w.Fair), nil
	case "label":
		return starlarkListFromDirectives(w.Label), nil
	case "machine_type":
		return starlarkListFromDirectives(w.MachineType), nil
	case "max_submit_await":
		return starlarkListFromDirectives(w.MaxSubmitAwait), nil
	case "max_errors":
		return starlarkListFromDirectives(w.MaxErrors), nil
	case "max_forks":
		return starlarkListFromDirectives(w.MaxForks), nil
	case "max_retries":
		return starlarkListFromDirectives(w.MaxRetries), nil
	case "memory":
		return starlarkListFromDirectives(w.Memory), nil
	case "module":
		return starlarkListFromDirectives(w.Module), nil
	case "penv":
		return starlarkListFromDirectives(w.Penv), nil
	case "pod":
		return starlarkListFromDirectives(w.Pod), nil
	case "publish_dir":
		return starlarkListFromDirectives(w.PublishDir), nil
	case "queue":
		return starlarkListFromDirectives(w.Queue), nil
	case "resource_labels":
		return starlarkListFromDirectives(w.ResourceLabels), nil
	case "resource_limits":
		return starlarkListFromDirectives(w.ResourceLimits), nil
	case "scratch":
		return starlarkListFromDirectives(w.Scratch), nil
	case "shell":
		return starlarkListFromDirectives(w.Shell), nil
	case "spack":
		return starlarkListFromDirectives(w.Spack), nil
	case "stage_in_mode":
		return starlarkListFromDirectives(w.StageInMode), nil
	case "stage_out_mode":
		return starlarkListFromDirectives(w.StageOutMode), nil
	case "store_dir":
		return starlarkListFromDirectives(w.StoreDir), nil
	case "tag":
		return starlarkListFromDirectives(w.Tag), nil
	case "time":
		return starlarkListFromDirectives(w.Time), nil
	case "dynamic":
		return starlarkListFromDirectives(w.Dynamic), nil
	case "unknown":
		return starlarkListFromDirectives(w.Unknown), nil
	default:
		return nil, fmt.Errorf("directives has no attribute %q", name)
	}
}

func (w *StarlarkProcessDirectivesWrapper) AttrNames() []string {
	return []string{
		"accelerator",
		"after_script",
		"arch",
		"array",
		"before_script",
		"cache",
		"cluster_options",
		"conda",
		"container",
		"container_options",
		"cpus",
		"debug",
		"disk",
		"echo",
		"error_strategy",
		"executor",
		"ext",
		"fair",
		"label",
		"machine_type",
		"max_submit_await",
		"max_errors",
		"max_forks",
		"max_retries",
		"memory",
		"module",
		"penv",
		"pod",
		"publish_dir",
		"queue",
		"resource_labels",
		"resource_limits",
		"scratch",
		"shell",
		"spack",
		"stage_in_mode",
		"stage_out_mode",
		"store_dir",
		"tag",
		"time",
		"dynamic",
		"unknown",
	}
}
