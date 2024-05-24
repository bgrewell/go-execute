package execute

import (
	"context"
	"io"
	"time"
)

type Executor interface {
	Execute(command string) (combined string, err error)
	ExecuteSeparate(command string) (stdout string, stderr string, err error)
	ExecuteAsync(command string) (result *ExecutionResult, err error)
	ExecuteAsyncWithInput(command string, stdin io.ReadCloser) (result *ExecutionResult, err error)
	ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error)
	ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error)
	ExecuteAsyncWithTimeout(command string, timeout time.Duration) (result *ExecutionResult, err error)
	ExecuteTTY(command string) error
}

// ExecutionResult holds the necessary structures for interaction with the process.
type ExecutionResult struct {
	Stdout   io.Reader
	Stderr   io.Reader
	Finished <-chan error
	Ctx      context.Context
}
