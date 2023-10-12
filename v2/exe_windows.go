package v2

import (
	"context"
	"github.com/BGrewell/go-execute/internal/utilities"
	"io"
	"os"
	"os/exec"
	"time"
)

const (
	SE_ASSIGNPRIMARYTOKEN_NAME = "SeAssignPrimaryTokenPrivilege"
	SE_INCREASE_QUOTA_NAME     = "SeIncreaseQuotaPrivilege"
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
		hasPrivileges, err := hasRequiredPrivileges()
		if err != nil {
			cancel()
			return nil, nil, err
		}

		// Check if the target user exists on the system
		sid, err := lookupAccount(e.User)
		if err != nil {
			cancel()
			return nil, nil, err
		}
	}

	return exe, cancel, nil
}

func hasRequiredPrivileges() (bool, error) {
	var hToken syscall.Token
	err := syscall.OpenProcessToken(syscall.CurrentProcess(), syscall.TOKEN_QUERY, &hToken)
	if err != nil {
		return false, err
	}
	defer syscall.CloseHandle(hToken)

	tokenPrivs, err := hToken.GetTokenPrivileges()
	if err != nil {
		return false, err
	}

	hasAssignPrimaryToken := false
	hasIncreaseQuota := false
	for _, priv := range tokenPrivs {
		name, err := priv.Name()
		if err != nil {
			continue
		}
		if name == SE_ASSIGNPRIMARYTOKEN_NAME {
			hasAssignPrimaryToken = true
		}
		if name == SE_INCREASE_QUOTA_NAME {
			hasIncreaseQuota = true
		}
	}

	return hasAssignPrimaryToken && hasIncreaseQuota, nil
}

func lookupAccount(username string) (*syscall.SID, error) {
	// Use the LookupAccountName syscall to verify the user exists
	var sid *syscall.SID
	var domain *uint16
	var size uint32
	var peUse uint32
	err := syscall.LookupAccountName(nil, syscall.StringToUTF16Ptr(username), sid, &size, &domain, &size, &peUse)
	return sid, err
}
