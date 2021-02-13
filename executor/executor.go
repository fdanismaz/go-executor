package executor

type Task interface {
	Run() error
	GetName() string
}

type Executor interface {
	Submit(t Task)
	Execute() error
	IsFinished() bool
}
