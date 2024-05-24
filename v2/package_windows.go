package execute

import (
	"io"
	"os"
	"time"
)

var defaultExecutor = NewExecutorWithEnv(os.Environ())

func Execute(command string) (string, error) {
	return defaultExecutor.Execute(command)
}

func ExecuteWithTimeout(command string, timeout time.Duration) (string, error) {
	return defaultExecutor.ExecuteWithTimeout(command, timeout)
}

func ExecuteSeparate(command string) (string, string, error) {
	return defaultExecutor.ExecuteSeparate(command)
}

func ExecuteSeparateWithTimeout(command string, timeout time.Duration) (string, string, error) {
	return defaultExecutor.ExecuteSeparateWithTimeout(command, timeout)
}

func ExecuteStream(command string) (io.ReadCloser, io.ReadCloser, error) {
	return defaultExecutor.ExecuteStream(command)
}

func ExecuteStreamWithTimeout(command string, timeout time.Duration) (io.ReadCloser, io.ReadCloser, error) {
	return defaultExecutor.ExecuteStreamWithTimeout(command, timeout)
}

func ExecuteStreamWithInput(command string, stdin io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	return defaultExecutor.ExecuteStreamWithInput(command, stdin)
}

func ExecuteTTY(command string) error {
	return defaultExecutor.ExecuteTTY(command)
}
