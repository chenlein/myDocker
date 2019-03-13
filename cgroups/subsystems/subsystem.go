package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	Name() string
	Set(path string, resource *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var (
	SubsystemInstances = []Subsystem{
		&MemorySubsystem{},
	}
)
