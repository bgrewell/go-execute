package execute

import (
	"context"
	"errors"
	"fmt"
	"github.com/bgrewell/go-execute/v2/internal/utilities"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"os/exec"
	"syscall"
)

func NewExecutor() Executor {
	return NewExecutorAsUser("", os.Environ())
}

func NewExecutorWithEnv(env []string) Executor {
	return NewExecutorAsUser("", env)
}

func NewExecutorAsUser(user string, env []string) Executor {
	return &WindowsExecutor{
		Environment: env,
		User:        user,
	}
}

type WindowsExecutor struct {
	BaseExecutor
	Environment []string
	User        string
}

//func (e WindowsExecutor) Execute(command string) (combined string, err error) {
//	return e.ExecuteWithTimeout(command, 0)
//}
//
//func (e WindowsExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
//	return e.ExecuteSeparateWithTimeout(command, 0)
//}
//
//func (e WindowsExecutor) ExecuteAsync(command string) (*ExecutionResult, error) {
//	return e.ExecuteAsyncWithTimeout(command, 0)
//}
//
//func (e WindowsExecutor) ExecuteAsyncWithInput(command string, stdin io.ReadCloser) (*ExecutionResult, error) {
//	return e.executeAsync(command, stdin, 0)
//}
//
//func (e WindowsExecutor) ExecuteAsyncWithTimeout(command string, timeout time.Duration) (*ExecutionResult, error) {
//	return e.executeAsync(command, nil, timeout)
//}
//
//func (e WindowsExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
//	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
//	return sout + serr, err
//}
//
//func (e WindowsExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
//	sout, serr, err := e.execute(command, nil, timeout)
//	if err != nil {
//		return "", "", err
//	}
//
//	outBytes, err := io.ReadAll(sout)
//	if err != nil {
//		return "", "", err
//	}
//	errBytes, err := io.ReadAll(serr)
//	if err != nil {
//		return "", "", err
//	}
//
//	return string(outBytes), string(errBytes), nil
//}
//
//func (e WindowsExecutor) ExecuteTTY(command string) error {
//	exe, _, cancel, err := e.prepareCommand(command, os.Stdin, 0)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if cancel != nil {
//			cancel()
//		}
//	}()
//
//	exe.Stdout = os.Stdout
//	exe.Stderr = os.Stderr
//
//	err = exe.Start()
//	if err != nil {
//		return err
//	}
//
//	return exe.Wait()
//}
//
//func (e WindowsExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
//	exe, ctx, cancel, err := e.prepareCommand(command, stdin, timeout)
//	if err != nil {
//		return nil, nil, err
//	}
//	defer func() {
//		if cancel != nil {
//			cancel()
//		}
//	}()
//
//	stdout, err = exe.StdoutPipe()
//	if err != nil {
//		return nil, nil, err
//	}
//
//	stderr, err = exe.StderrPipe()
//	if err != nil {
//		return nil, nil, err
//	}
//
//	var stdoutBuf, stderrBuf bytes.Buffer
//	stdoutDone := make(chan struct{})
//	stderrDone := make(chan struct{})
//
//	go utilities.CopyAndClose(stdoutDone, &stdoutBuf, stdout)
//	go utilities.CopyAndClose(stderrDone, &stderrBuf, stderr)
//
//	err = exe.Start()
//	if err != nil {
//		return nil, nil, err
//	}
//
//	// Wait for the command to complete or for the timeout
//	done := make(chan error, 1)
//	go func() {
//		done <- exe.Wait()
//	}()
//
//	// Wait for completion or timeout
//	select {
//	case <-ctx.Done():
//		return nil, nil, ctx.Err()
//	case err := <-done:
//		return io.NopCloser(bytes.NewReader(stdoutBuf.Bytes())), io.NopCloser(bytes.NewReader(stderrBuf.Bytes())), err
//	}
//}
//
//func (e WindowsExecutor) executeAsync(command string, stdin io.ReadCloser, timeout time.Duration) (*ExecutionResult, error) {
//	exe, ctx, cancel, err := e.prepareCommand(command, stdin, timeout)
//	if err != nil {
//		return nil, err
//	}
//
//	// Setting up stdout and stderr
//	stdoutPipe, err := exe.StdoutPipe()
//	if err != nil {
//		cancel()
//		return nil, err
//	}
//	stderrPipe, err := exe.StderrPipe()
//	if err != nil {
//		cancel()
//		return nil, err
//	}
//
//	// Buffering stdout and stderr
//	var stdoutBuf, stderrBuf bytes.Buffer
//	stdoutDone := make(chan struct{})
//	stderrDone := make(chan struct{})
//
//	go utilities.CopyAndClose(stdoutDone, &stdoutBuf, stdoutPipe)
//	go utilities.CopyAndClose(stderrDone, &stderrBuf, stderrPipe)
//
//	// Starting the command asynchronously
//	err = exe.Start()
//	if err != nil {
//		if cancel != nil {
//			cancel()
//		}
//		return nil, err
//	}
//
//	finished := make(chan error)
//	go func() {
//		defer close(finished)
//		finished <- exe.Wait()
//		if cancel != nil {
//			cancel()
//		}
//	}()
//
//	return &ExecutionResult{
//		Stdout:   &stdoutBuf,
//		Stderr:   &stderrBuf,
//		Finished: finished,
//		Ctx:      ctx,
//	}, nil
//}
//
//func (e WindowsExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
//	ctx := context.Background()
//	var cancel context.CancelFunc
//
//	if timeout != 0 {
//		ctx, cancel = context.WithTimeout(ctx, timeout)
//	}
//
//	cmdParts, err := utilities.Fields(command) // Assuming utilities.Fields breaks the command string into parts
//	if err != nil {
//		return nil, ctx, cancel, err
//	}
//
//	binary, err := exec.LookPath(cmdParts[0])
//	if err != nil {
//		return nil, ctx, cancel, err
//	}
//
//	exe := exec.CommandContext(ctx, binary, cmdParts[1:]...)
//	exe.Stdin = stdin
//	exe.Env = e.Environment
//
//	if e.User != "" {
//		err = e.configureUser(ctx, cancel, exe)
//		if err != nil {
//			return nil, ctx, cancel, err
//		}
//	}
//
//	return exe, ctx, cancel, nil
//}

func (e WindowsExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	// Check if the current process has the required privileges
	isAdmin := utilities.RunningAsAdmin()
	fmt.Printf("IsAdmin: %v\n", isAdmin)

	// Try to find a process running as the target user
	pid := int32(0)
	processes, _ := process.Processes()
	for _, process := range processes {
		username, _ := process.Username()
		if username == e.User {
			pid = process.Pid
			break
		}
	}
	if pid == 0 {
		return errors.New("unable to find process running as target user")
	}

	token, err := utilities.GetTokenFromPid(pid)
	if err != nil {
		return err
	}

	exe.SysProcAttr = &syscall.SysProcAttr{
		Token: token,
	}

	return nil
}

//func (e WindowsExecutor) ExecutePowershell(command string) (combined string, err error) {
//    return e.ExecuteWithTimeout(e.encodePowershellCommand(command), 0)
//}
//
//func (e WindowsExecutor) ExecutePowershellSeparate(command string) (stdout string, stderr string, err error) {
//    return e.ExecuteSeparateWithTimeout(e.encodePowershellCommand(command), 0)
//}
//
//func (e WindowsExecutor) ExecutePowershellStream(command string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
//    return e.ExecuteStreamWithTimeout(e.encodePowershellCommand(command), 0)
//}
//
//func (e WindowsExecutor) ExecutePowershellStreamWithInput(command string, stdin io.ReadCloser) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
//    return e.execute(e.encodePowershellCommand(command), stdin, 0)
//}
//
//func (e WindowsExecutor) ExecutePowershellWithTimeout(command string, timeout time.Duration) (combined string, err error) {
//    sout, serr, err := e.ExecuteSeparateWithTimeout(e.encodePowershellCommand(command), timeout)
//    return sout + serr, err
//}
//
//func (e WindowsExecutor) ExecutePowershellSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
//    sout, serr, err := e.execute(e.encodePowershellCommand(command), nil, timeout)
//    if err != nil {
//        return "", "", err
//    }
//
//    outBytes, _ := io.ReadAll(sout)
//    errBytes, _ := io.ReadAll(serr)
//
//    return string(outBytes), string(errBytes), nil
//}
//
//func (e WindowsExecutor) ExecutePowershellStreamWithTimeout(command string, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
//    return e.execute(e.encodePowershellCommand(command), nil, timeout)
//}
//
//func (e WindowsExecutor) ExecutePowershellTTY(command string) error {
//    exe, cancel, err := e.prepareCommand(e.encodePowershellCommand(command), os.Stdin, 0)
//    if err != nil {
//        return err
//    }
//    defer func() {
//        if cancel != nil {
//            cancel()
//        }
//    }()
//
//    exe.Stdout = os.Stdout
//    exe.Stderr = os.Stderr
//
//    err = exe.Start()
//    if err != nil {
//        return err
//    }
//
//    return exe.Wait()
//}
//
//func (e WindowsExecutor) encodePowershellCommand(command string) string {
//    encCommand := utilities.ConvertToUTF16LEBase64String(command)
//    return fmt.Sprintf("powershell.exe -NoProfile -enc %s", encCommand)
//}
