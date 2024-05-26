package execute

import (
	"bytes"
	"context"
	"errors"
	"github.com/bgrewell/go-execute/v2/internal/utilities"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Executor is the interface that wraps the basic Execute functions.
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

// BaseExecutor is the base implementation of the Executor interface. It implements all the code that is shared between
// the platform-specific executors.
type BaseExecutor struct {
	Environment []string
	User        string
}

// Execute is the base implementation of the Execute function which executes a command and returns the combined output.
func (e BaseExecutor) Execute(command string) (combined string, err error) {
	return e.ExecuteWithTimeout(command, 0)
}

// ExecuteSeparate is the base implementation of the ExecuteSeparate function which executes a command and returns the stdout and stderr separately.
func (e BaseExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	return e.ExecuteSeparateWithTimeout(command, 0)
}

// ExecuteAsync is the base implementation of the ExecuteAsync function which executes a command asynchronously.
func (e BaseExecutor) ExecuteAsync(command string) (*ExecutionResult, error) {
	return e.ExecuteAsyncWithTimeout(command, 0)
}

// ExecuteAsyncWithInput is the base implementation of the ExecuteAsyncWithInput function which executes a command asynchronously with input.
func (e BaseExecutor) ExecuteAsyncWithInput(command string, stdin io.ReadCloser) (*ExecutionResult, error) {
	return e.executeAsync(command, stdin, 0)
}

// ExecuteAsyncWithTimeout is the base implementation of the ExecuteAsyncWithTimeout function which executes a command asynchronously with a timeout.
func (e BaseExecutor) ExecuteAsyncWithTimeout(command string, timeout time.Duration) (*ExecutionResult, error) {
	return e.executeAsync(command, nil, timeout)
}

// ExecuteWithTimeout is the base implementation of the ExecuteWithTimeout function which executes a command with a timeout.
func (e BaseExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
	return sout + serr, err
}

// ExecuteSeparateWithTimeout is the base implementation of the ExecuteSeparateWithTimeout function which executes a command and returns the stdout and stderr separately with a timeout.
func (e BaseExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
	sout, serr, err := e.execute(command, nil, timeout)
	if err != nil {
		return "", "", err
	}

	outBytes, err := io.ReadAll(sout)
	if err != nil {
		return "", "", err
	}
	errBytes, err := io.ReadAll(serr)
	if err != nil {
		return "", "", err
	}

	return string(outBytes), string(errBytes), nil
}

// ExecuteTTY is the base implementation of the ExecuteTTY function which executes a command with a TTY.
func (e BaseExecutor) ExecuteTTY(command string) error {
	exe, _, cancel, err := e.prepareCommand(command, os.Stdin, 0)
	if err != nil {
		return err
	}
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr

	err = exe.Start()
	if err != nil {
		return err
	}

	return exe.Wait()
}

// execute is the base implementation of the execute function which executes a command and returns the stdout and stderr.
func (e BaseExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (io.ReadCloser, io.ReadCloser, error) {
	execResult, err := e.executeAsync(command, stdin, timeout)
	if err != nil {
		logger.Error("failed to execute command", "error", err)
		return nil, nil, err
	}

	// Wait for completion or timeout using the context from execResult
	select {
	case err := <-execResult.Finished:
		logger.Trace("command execution finished")
		return io.NopCloser(bytes.NewReader(execResult.Stdout.(*bytes.Buffer).Bytes())),
			io.NopCloser(bytes.NewReader(execResult.Stderr.(*bytes.Buffer).Bytes())),
			err
	case <-execResult.Ctx.Done():
		logger.Error("command execution timed out", "error", execResult.Ctx.Err())
		return nil, nil, execResult.Ctx.Err()
	}
}

// executeAsync is the base implementation of the executeAsync function which executes a command asynchronously.
func (e BaseExecutor) executeAsync(command string, stdin io.ReadCloser, timeout time.Duration) (*ExecutionResult, error) {
	exe, ctx, cancel, err := e.prepareCommand(command, stdin, timeout)
	if err != nil {
		return nil, err
	}

	// Setting up stdout and stderr
	stdoutPipe, err := exe.StdoutPipe()
	if err != nil {
		logger.Error("failed to get stdout pipe", "error", err)
		cancel()
		return nil, err
	}
	stderrPipe, err := exe.StderrPipe()
	if err != nil {
		logger.Error("failed to get stderr pipe", "error", err)
		cancel()
		return nil, err
	}

	// Buffering stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})
	logger.Trace("captured stdout and stderr pipes")

	var ready sync.WaitGroup
	ready.Add(2)
	go utilities.CopyAndClose(stdoutDone, &stdoutBuf, stdoutPipe, &ready)
	go utilities.CopyAndClose(stderrDone, &stderrBuf, stderrPipe, &ready)
	logger.Trace("waiting for pipe readers to start")
	ready.Wait()
	logger.Trace("pipe readers have started")

	// Starting the command asynchronously
	err = exe.Start()
	if err != nil {
		logger.Error("failed to start command", "error", err)
		if cancel != nil {
			cancel()
		}
		return nil, err
	}
	logger.Trace("started command asynchronously")

	finished := make(chan error)
	go func() {
		defer close(finished)
		exitErr := exe.Wait()
		finished <- exitErr
		logger.Trace("command finished executing", "exit", exitErr)
		if cancel != nil {
			cancel()
		}
	}()

	logger.Trace("returning ExecutionResults object")
	return &ExecutionResult{
		Stdout:   &stdoutBuf,
		Stderr:   &stderrBuf,
		Finished: finished,
		Ctx:      ctx,
	}, nil
}

// prepareCommand is the base implementation of the prepareCommand function which prepares the command for execution.
func (e BaseExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout != 0 {
		logger.Trace("configuring command timeout", "timeout", timeout)
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command) // Assuming utilities.Fields breaks the command string into parts
	if err != nil {
		logger.Error("failed to get command parts", "error", err)
		return nil, ctx, cancel, err
	}
	logger.Trace("split command into the parts", "cmdParts", cmdParts)

	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		logger.Error("failed to find binary path", "error", err)
		return nil, ctx, cancel, err
	}
	logger.Trace("binary found", "binary", binary)

	exe := exec.CommandContext(ctx, binary, cmdParts[1:]...)
	exe.Stdin = stdin
	exe.Env = e.Environment
	logger.Trace("command context set", "environment", exe.Env)

	if e.User != "" {
		err := e.configureUser(ctx, cancel, exe)
		if err != nil {
			logger.Error("failed to configure command user", "error", err)
			return exe, ctx, cancel, err
		}
		logger.Trace("configured user for execution", "user", e.User)
	}

	return exe, ctx, cancel, nil
}

// configureUser is the base implementation of the configureUser function which must be overridden by the platform-specific executor.
func (e BaseExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	return errors.New("this method must be implemented by the platform-specific executor")
}

// Struct and methods to allow basic execution without needing to instantiate a new Executor.
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
