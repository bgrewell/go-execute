package execute

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bgrewell/go-execute/v2/internal/utilities"
	"github.com/shirou/gopsutil/v3/process"
	"io"
	"os"
	"os/exec"
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
	return e.ExecuteWithTimeout(command, 0)
}

func (e WindowsExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	return e.ExecuteSeparateWithTimeout(command, 0)
}

func (e WindowsExecutor) ExecuteStream(command string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	return e.ExecuteStreamWithTimeout(command, 0)
}

func (e WindowsExecutor) ExecuteStreamWithInput(command string, stdin io.ReadCloser) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	return e.execute(command, stdin, 0)
}

func (e WindowsExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
	return sout + serr, err
}

func (e WindowsExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
	sout, serr, err := e.execute(command, nil, timeout)
	if err != nil {
		return "", "", err
	}

	outBytes, _ := io.ReadAll(sout)
	errBytes, _ := io.ReadAll(serr)

	return string(outBytes), string(errBytes), nil
}

func (e WindowsExecutor) ExecuteStreamWithTimeout(command string, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	return e.execute(command, nil, timeout)
}

func (e WindowsExecutor) ExecuteTTY(command string) error {
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

func (e WindowsExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
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

	go utilities.CopyAndClose(stdoutDone, &stdoutBuf, stdout)
	go utilities.CopyAndClose(stderrDone, &stderrBuf, stderr)

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

func (e WindowsExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
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
			return nil, ctx, cancel, errors.New("unable to find process running as target user")
		}

		token, err := utilities.GetTokenFromPid(pid)
		if err != nil {
			return nil, ctx, cancel, err
		}

		exe.SysProcAttr = &syscall.SysProcAttr{
			Token: token,
		}
	}

	return exe, ctx, cancel, nil
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
