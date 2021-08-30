package process

var (
	processManager Process
)

// Process provides interfaces to query process.
type Process interface {
	IsRunning(name string) (bool, error)
}

// NewProcessManager creates process instance.
func NewProcessManager() Process {
	return &process{}
}
