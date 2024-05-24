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

type PowerShellExecutor struct {
	Environment []string
	User        string
	Session     *exec.Cmd
}

func NewPowerShellExecutor() PowerShellExecutor {
	return NewPowerShellExecutorAsUser("", os.Environ())
}

func NewPowerShellExecutorWithEnv(env []string) PowerShellExecutor {
	return NewPowerShellExecutorAsUser("", env)
}

func NewPowerShellExecutorAsUser(user string, env []string) PowerShellExecutor {
	return PowerShellExecutor{
		Environment: env,
		User:        user,
	}
}

func (e PowerShellExecutor) Execute(command string) (string, error) {
	return e.ExecuteWithTimeout(command, 0)
}

func (e PowerShellExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (string, error) {
	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
	return sout + serr, err
}

func (e PowerShellExecutor) ExecuteSeparate(command string) (string, string, error) {
	return e.ExecuteSeparateWithTimeout(command, 0)
}

func (e PowerShellExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (string, string, error) {
	sout, serr, err := e.execute(command, nil, timeout)
	if err != nil {
		return "", "", err
	}

	outBytes, _ := io.ReadAll(sout)
	errBytes, _ := io.ReadAll(serr)

	return string(outBytes), string(errBytes), nil
}

func (e PowerShellExecutor) ExecuteScript(script string, parameters map[string]string) (string, error) {
	command := "powershell -NoProfile -NonInteractive -Command " + script
	for key, value := range parameters {
		command += fmt.Sprintf(" -%s '%s'", key, value)
	}
	return e.Execute(command)
}

func (e PowerShellExecutor) StartSession() error {
	e.Session = exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "-")
	e.Session.Env = e.Environment
	if e.User != "" {
		token, err := utilities.GetTokenForUser(e.User)
		if err != nil {
			return err
		}
		e.Session.SysProcAttr = &syscall.SysProcAttr{
			Token: token,
		}
	}

	stdin, err := e.Session.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()

	return e.Session.Start()
}

func (e PowerShellExecutor) ExecuteInSession(command string) (string, error) {
	if e.Session == nil {
		return "", errors.New("PowerShell session not started")
	}

	stdin, err := e.Session.StdinPipe()
	if err != nil {
		return "", err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	e.Session.Stdout = &stdoutBuf
	e.Session.Stderr = &stderrBuf

	_, err = stdin.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	err = e.Session.Wait()
	if err != nil {
		return "", err
	}

	return stdoutBuf.String() + stderrBuf.String(), nil
}

func (e PowerShellExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
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

func (e PowerShellExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command)
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
