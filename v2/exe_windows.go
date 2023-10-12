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

func (e WindowsExecutor) ExecuteTTY(command string, timeout time.Duration) error {
	exe, cancel, err := e.prepareCommandWindows(command, os.Stdin, timeout)
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

func (e WindowsExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	exe, cancel, err := e.prepareCommand(command, stdin, timeout)
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

	err = exe.Start()
	if err != nil {
		return nil, nil, err
	}

	go func() {
		_ = exe.Wait()
	}()

	return stdout, stderr, nil
}

func (e WindowsExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command) // Assuming utilities.Fields breaks the command string into parts
	if err != nil {
		return nil, cancel, err
	}

	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return nil, cancel, err
	}

	exe := exec.CommandContext(ctx, binary, cmdParts[1:]...)
	exe.Stdin = stdin
	exe.Env = e.Environment

	if e.User != "" {
		// This part will be more complex on Windows and may involve direct syscall
		// operations or third-party libraries to support user-based execution.
		// For now, let's leave it as a placeholder.
	}

	return exe, cancel, nil
}
