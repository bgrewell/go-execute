package v2

import (
	"bytes"
	"context"
	"github.com/BGrewell/go-execute/internal/utilities"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func NewExecutor(env []string) Executor {
	return NewExecutorAsUser("", os.Environ())
}

func NewExecutorWithEnv(env []string) Executor {
	return NewExecutorAsUser("", env)
}

func NewExecutorAsUser(user string, env []string) Executor {
	return &DarwinExecutor{
		Environment: env,
		User:        user,
	}
}

type DarwinExecutor struct {
	Environment []string
	User        string
}

func (e DarwinExecutor) Execute(command string) (combined string, err error) {
	return e.ExecuteWithTimeout(command, 0)
}

func (e DarwinExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	return e.ExecuteSeparateWithTimeout(command, 0)
}

func (e DarwinExecutor) ExecuteAsync(command string) (*ExecutionResult, error) {
	return e.ExecuteAsyncWithTimeout(command, 0)
}

func (e DarwinExecutor) ExecuteAsyncWithInput(command string, stdin io.ReadCloser) (*ExecutionResult, error) {
	return e.executeAsync(command, stdin, 0)
}

func (e DarwinExecutor) ExecuteAsyncWithTimeout(command string, timeout time.Duration) (*ExecutionResult, error) {
	return e.executeAsync(command, nil, timeout)
}

func (e DarwinExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
	return sout + serr, err
}

func (e DarwinExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
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

func (e DarwinExecutor) ExecuteTTY(command string) error {
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

func (e DarwinExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (io.ReadCloser, io.ReadCloser, error) {
	execResult, err := e.executeAsync(command, stdin, timeout)
	if err != nil {
		return nil, nil, err
	}

	// Wait for completion or timeout using the context from execResult
	select {
	case err := <-execResult.Finished:
		return io.NopCloser(bytes.NewReader(execResult.Stdout.(*bytes.Buffer).Bytes())),
			io.NopCloser(bytes.NewReader(execResult.Stderr.(*bytes.Buffer).Bytes())),
			err
	case <-execResult.Ctx.Done():
		return nil, nil, execResult.Ctx.Err()
	}
}

// executeAsync starts the command asynchronously and returns access to stdout, stderr, and a completion channel.
func (e DarwinExecutor) executeAsync(command string, stdin io.ReadCloser, timeout time.Duration) (*ExecutionResult, error) {
	exe, ctx, cancel, err := e.prepareCommand(command, stdin, timeout)
	if err != nil {
		return nil, err
	}

	// Setting up stdout and stderr
	stdoutPipe, err := exe.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stderrPipe, err := exe.StderrPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	// Buffering stdout and stderr
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go copyAndClose(stdoutDone, &stdoutBuf, stdoutPipe)
	go copyAndClose(stderrDone, &stderrBuf, stderrPipe)

	// Starting the command asynchronously
	err = exe.Start()
	if err != nil {
		if cancel != nil {
			cancel()
		}
		return nil, err
	}

	finished := make(chan error)
	go func() {
		defer close(finished)
		finished <- exe.Wait()
		if cancel != nil {
			cancel()
		}
	}()

	return &ExecutionResult{
		Stdout:   &stdoutBuf,
		Stderr:   &stderrBuf,
		Finished: finished,
		Ctx:      ctx,
	}, nil
}

func (e DarwinExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command) // Assuming utilities.Fields breaks the command string into parts
	if err != nil {
		return nil, ctx, cancel, err
	}

	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return nil, ctx, cancel, err
	}

	exe := exec.CommandContext(ctx, binary, cmdParts[1:]...)
	exe.Stdin = stdin
	exe.Env = e.Environment

	if e.User != "" {
		u, err := user.Lookup(e.User)
		if err != nil {
			return nil, ctx, cancel, err
		}

		uid, err := strconv.Atoi(u.Uid)
		if err != nil {
			return nil, ctx, cancel, err
		}

		gid, err := strconv.Atoi(u.Gid)
		if err != nil {
			return nil, ctx, cancel, err
		}

		exe.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
		}
	}

	return exe, ctx, cancel, nil
}

func copyAndClose(done chan struct{}, buf *bytes.Buffer, r io.ReadCloser) {
	defer close(done)
	defer r.Close()
	io.Copy(buf, r)
}
