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

func NewExecutor() Executor {
	return NewExecutorAsUser("", os.Environ())
}

func NewExecutorWithEnv(env []string) Executor {
	return NewExecutorAsUser("", env)
}

func NewExecutorAsUser(user string, env []string) Executor {
	return &LinuxExecutor{
		Environment: env,
		User:        user,
	}
}

type LinuxExecutor struct {
	Environment []string
	User        string
}

func (e LinuxExecutor) Execute(command string) (combined string, err error) {
	return e.ExecuteWithTimeout(command, 0)
}

func (e LinuxExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	return e.ExecuteSeparateWithTimeout(command, 0)
}

func (e LinuxExecutor) ExecuteStream(command string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	return e.ExecuteStreamWithTimeout(command, 0)
}

func (e LinuxExecutor) ExecuteStreamWithInput(command string, stdin io.ReadCloser) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	return e.execute(command, stdin, 0)
}

func (e LinuxExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
	return sout + serr, err
}

func (e LinuxExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
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

func (e LinuxExecutor) ExecuteStreamWithTimeout(command string, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	return e.execute(command, nil, timeout)
}

func (e LinuxExecutor) ExecuteTTY(command string) error {
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

func (e LinuxExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	exe, ctx, cancel, err := e.prepareCommand(command, stdin, timeout)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	stdout, err = exe.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err = exe.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go copyAndClose(stdoutDone, &stdoutBuf, stdout)
	go copyAndClose(stderrDone, &stderrBuf, stderr)

	err = exe.Start()
	if err != nil {
		return nil, nil, err
	}

	// Wait for the command to complete or for the timeout
	done := make(chan error, 1)
	go func() {
		done <- exe.Wait()
	}()

	// Wait for completion or timeout
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case err := <-done:
		return io.NopCloser(bytes.NewReader(stdoutBuf.Bytes())), io.NopCloser(bytes.NewReader(stderrBuf.Bytes())), err
	}
}

func (e LinuxExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
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
