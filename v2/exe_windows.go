package execute

import (
	"context"
	"errors"
	"fmt"
	"github.com/bgrewell/go-execute/v2/internal/utilities"
	"github.com/shirou/gopsutil/v3/process"
	"os/exec"
	"syscall"
)

// NewExecutor creates a new Executor.
func NewExecutor(options ...Option) Executor {
	e := &WindowsExecutor{}
	for _, option := range options {
		option(e)
	}
	return e
}

// WindowsExecutor is an Executor implementation for Windows systems.
type WindowsExecutor struct {
	BaseExecutor
}

// configureUser sets the user and group for the command to be executed.
func (e WindowsExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	// Check if the current process has the required privileges
	isAdmin := utilities.RunningAsAdmin()
	fmt.Printf("IsAdmin: %v\n", isAdmin)

	// Try to find a process running as the target user
	pid := int32(0)
	processes, _ := process.Processes()
	for _, process := range processes {
		username, _ := process.Username()
		if username == e.user {
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
