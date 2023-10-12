package v2

import (
	"io"
	"time"
)

func NewExecutor(env []string) Executor {
	return NewExecutorAsUser("", env)
}

func NewExecutorAsUser(user string, env []string) Executor {
	return &WindowsExecutor{
		Environment: env,
		User:        user,
	}
}

type WindowsExecutor struct {
	Environment []string
	User        string
}

func (e WindowsExecutor) Execute(command string) (combined string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteStream(command string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteStreamWithInput(command string, stdin io.WriteCloser) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteStreamWithTimeout(command string, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}

func (e WindowsExecutor) ExecuteTTY(command string) error {
	panic("implement me")
}
