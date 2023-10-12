package v2

import (
	"context"
	"errors"
	"fmt"
	"github.com/BGrewell/go-execute/internal/utilities"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows"
	"io"
	"os"
	"os/exec"
	"syscall"
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
	exe, cancel, err := e.prepareCommand(command, os.Stdin, 0)
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

// execute contains Windows specific execution code which is called from the various public methods
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

// prepareCommand contains Windows specific command execution prep which is called from the various public methods
func (e WindowsExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	// Helper function
	cf := func(err error) (*exec.Cmd, context.CancelFunc, error) {
		if cancel != nil {
			cancel()
		}
		return nil, nil, err
	}

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command)
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
		// Check if the current process has the required privileges
		isAdmin := runningAsAdmin()
		if err != nil {
			return cf(err)
		}
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
			return cf(errors.New("unable to find process running as target user"))
		}

		token, err := getTokenFromPid(pid)
		if err != nil {
			return cf(err)
		}

		exe.SysProcAttr = &syscall.SysProcAttr{
			Token: token,
		}
	}

	return exe, cancel, nil
}

func hasRequiredPrivileges() (admin bool, elevated bool, err error) {
	var sid *windows.SID

	err = windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false, false, err
	}
	defer windows.FreeSid(sid)

	// Get the token for the active thread
	token := windows.Token(0)

	isAdmin, err := token.IsMember(sid)
	if err != nil {
		return false, false, err
	}

	return isAdmin, token.IsElevated(), nil
}

func runningAsAdmin() (isAdmin bool) {
	// Check to see if we are running with the right permissions
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func lookupAccount(username string) (*syscall.SID, error) {
	// Use the LookupAccountName syscall to verify the user exists
	var sid *syscall.SID
	var domain uint16
	var size uint32
	var peUse uint32
	err := syscall.LookupAccountName(nil, syscall.StringToUTF16Ptr(username), sid, &size, &domain, &size, &peUse)
	return sid, err
}

func getTokenFromPid(pid int32) (syscall.Token, error) {
	var err error
	var token syscall.Token

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		fmt.Println("Token Process", "err", err)
	}
	defer syscall.CloseHandle(handle)

	// Find process token via win32
	err = syscall.OpenProcessToken(handle, syscall.TOKEN_ALL_ACCESS, &token)

	if err != nil {
		fmt.Println("Open Token Process", "err", err)
	}
	return token, err
}
