package execute

import (
	"io"
	"os"
	"time"
)

var defaultExecutor = NewExecutorWithEnv(os.Environ())

func Execute(command string) (combined string, err error) {
	return defaultExecutor.Execute(command)
}

func ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	return defaultExecutor.ExecuteSeparate(command)
}

func ExecuteAsync(command string) (*ExecutionResult, error) {
	return defaultExecutor.ExecuteAsync(command)
}

func ExecuteAsyncWithTimeout(command string, timeout time.Duration) (*ExecutionResult, error) {
	return defaultExecutor.ExecuteAsyncWithTimeout(command, timeout)
}

func ExecuteAsyncWithInput(command string, stdin io.ReadCloser) (*ExecutionResult, error) {
	return defaultExecutor.ExecuteAsyncWithInput(command, stdin)
}

func ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	return defaultExecutor.ExecuteWithTimeout(command, timeout)
}

func ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
	return defaultExecutor.ExecuteSeparateWithTimeout(command, timeout)
}

func ExecuteTTY(command string) error {
	return defaultExecutor.ExecuteTTY(command)
}
