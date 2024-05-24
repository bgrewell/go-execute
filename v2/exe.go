package execute

import (
	"io"
	"time"
)

type Executor interface {
	Execute(command string) (combined string, err error)
	ExecuteSeparate(command string) (stdout string, stderr string, err error)
	ExecuteStream(command string) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	ExecuteStreamWithInput(command string, stdin io.ReadCloser) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error)
	ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error)
	ExecuteStreamWithTimeout(command string, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error)
	ExecuteTTY(command string) error
}
