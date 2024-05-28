package execute

import (
	"context"
	"errors"
	"fmt"
	"github.com/bgrewell/go-execute/v2/internal"
	"github.com/bgrewell/go-execute/v2/internal/utilities"
	"io"
	"os"
	"os/exec"
	"strings"
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
	ExecuteScriptFromString(scriptType ScriptType, script string, arguments []string, parameters map[string]string) (stdout string, stderr string, err error)
	ExecuteScriptFromStringWithTimeout(scriptType ScriptType, script string, arguments []string, parameters map[string]string, timeout time.Duration) (stdout string, stderr string, err error)
	ExecuteScriptFromFile(scriptType ScriptType, scriptPath string, arguments []string, parameters map[string]string) (stdout string, stderr string, err error)
	ExecuteScriptFromFileWithTimeout(scriptType ScriptType, scriptPath string, arguments []string, parameters map[string]string, timeout time.Duration) (stdout string, stderr string, err error)
	ExecuteTTY(command string) error
	SetEnvironment(env []string)
	Environment() []string
	SetUser(user string)
	User() string
	SetShell(shell string)
	ClearShell()
	Shell() string
	UsingShell() bool
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
	environment []string
	user        string
	shell       string
}

// SetEnvironment sets the environment for the executor.
func (e *BaseExecutor) SetEnvironment(env []string) {
	e.environment = env
}

// Environment returns the environment for the executor.
func (e *BaseExecutor) Environment() []string {
	return e.environment
}

// SetUser sets the user for the executor.
func (e *BaseExecutor) SetUser(user string) {
	e.user = user
}

// User returns the user for the executor.
func (e *BaseExecutor) User() string {
	return e.user
}

// SetShell sets the shell for the executor.
func (e *BaseExecutor) SetShell(shell string) {
	e.shell = shell
}

// ClearShell clears the shell for the executor.
func (e *BaseExecutor) ClearShell() {
	e.shell = ""
}

// Shell returns the shell for the executor.
func (e *BaseExecutor) Shell() string {
	return e.shell
}

// UsingShell returns whether the executor is using a shell.
func (e *BaseExecutor) UsingShell() bool {
	return e.shell != ""
}

// Execute is the base implementation of the Execute function which executes a command and returns the combined output.
func (e *BaseExecutor) Execute(command string) (combined string, err error) {
	return e.ExecuteWithTimeout(command, 0)
}

// ExecuteSeparate is the base implementation of the ExecuteSeparate function which executes a command and returns the stdout and stderr separately.
func (e *BaseExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	return e.ExecuteSeparateWithTimeout(command, 0)
}

// ExecuteAsync is the base implementation of the ExecuteAsync function which executes a command asynchronously.
func (e *BaseExecutor) ExecuteAsync(command string) (*ExecutionResult, error) {
	return e.ExecuteAsyncWithTimeout(command, 0)
}

// ExecuteAsyncWithInput is the base implementation of the ExecuteAsyncWithInput function which executes a command asynchronously with input.
func (e *BaseExecutor) ExecuteAsyncWithInput(command string, stdin io.ReadCloser) (*ExecutionResult, error) {
	return e.executeAsync(command, stdin, 0, false)
}

// ExecuteAsyncWithTimeout is the base implementation of the ExecuteAsyncWithTimeout function which executes a command asynchronously with a timeout.
func (e *BaseExecutor) ExecuteAsyncWithTimeout(command string, timeout time.Duration) (*ExecutionResult, error) {
	return e.executeAsync(command, nil, timeout, false)
}

// ExecuteWithTimeout is the base implementation of the ExecuteWithTimeout function which executes a command with a timeout.
func (e *BaseExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	sout, serr, err := e.ExecuteSeparateWithTimeout(command, timeout)
	return sout + serr, err
}

// ExecuteSeparateWithTimeout is the base implementation of the ExecuteSeparateWithTimeout function which executes a command and returns the stdout and stderr separately with a timeout.
func (e *BaseExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
	sout, serr, err := e.execute(command, nil, timeout, false)
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

// ExecuteScriptFromString is the base implementation of the ExecuteScriptFromString function which executes a script from a string.
func (e *BaseExecutor) ExecuteScriptFromString(scriptType ScriptType, script string, arguments []string, parameters map[string]string) (stdout string, stderr string, err error) {
	return e.ExecuteScriptFromStringWithTimeout(scriptType, script, arguments, parameters, 0)
}

// ExecuteScriptFromStringWithTimeout is the base implementation of the ExecuteScriptFromStringWithTimeout function which executes a script from a string with a timeout.
func (e *BaseExecutor) ExecuteScriptFromStringWithTimeout(scriptType ScriptType, script string, arguments []string, parameters map[string]string, timeout time.Duration) (stdout string, stderr string, err error) {
	tmpFile, err := e.writeTempScript(scriptType, script)
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tmpFile)

	return e.ExecuteScriptFromFileWithTimeout(scriptType, tmpFile, arguments, parameters, timeout)
}

// ExecuteScriptFromFile is the base implementation of the ExecuteScriptFromFile function which executes a script from a file.
func (e *BaseExecutor) ExecuteScriptFromFile(scriptType ScriptType, scriptPath string, arguments []string, parameters map[string]string) (stdout string, stderr string, err error) {
	return e.ExecuteScriptFromFileWithTimeout(scriptType, scriptPath, arguments, parameters, 0)
}

// ExecuteScriptFromFileWithTimeout is the base implementation of the ExecuteScriptFromFileWithTimeout function which executes a script from a file with a timeout.
func (e *BaseExecutor) ExecuteScriptFromFileWithTimeout(scriptType ScriptType, scriptPath string, arguments []string, parameters map[string]string, timeout time.Duration) (stdout string, stderr string, err error) {
	return e.executeScript(scriptType, scriptPath, arguments, parameters, timeout)
}

// ExecuteTTY is the base implementation of the ExecuteTTY function which executes a command with a TTY.
func (e *BaseExecutor) ExecuteTTY(command string) error {
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

// executeScript is the base implementation of the executeScript function which executes a script and returns the stdout and stderr.
func (e *BaseExecutor) executeScript(scriptType ScriptType, scriptPath string, arguments []string, parameters map[string]string, timeout time.Duration) (stdout string, stderr string, err error) {
	command, err := e.buildScriptCommand(scriptType, scriptPath, arguments, parameters)
	if err != nil {
		return "", "", err
	}

	sout, serr, err := e.execute(command, nil, timeout, true)

	outBytes, soerr := io.ReadAll(sout)
	if soerr != nil {
		if err != nil {
			return "", "", err
		}
		return "", "", soerr
	}
	errBytes, seerr := io.ReadAll(serr)
	if seerr != nil {
		if err != nil {
			return "", "", err
		}
		return "", "", seerr
	}

	return string(outBytes), string(errBytes), err
}

// execute is the base implementation of the execute function which executes a command and returns the stdout and stderr.
func (e *BaseExecutor) execute(command string, stdin io.ReadCloser, timeout time.Duration, script bool) (io.ReadCloser, io.ReadCloser, error) {
	execResult, err := e.executeAsync(command, stdin, timeout, script)
	if err != nil {
		logger.Error("failed to execute command", "error", err)
		return nil, nil, err
	}

	// Wait for completion or timeout using the context from execResult
	logger.Trace("waiting for command execution to finish")
	select {
	case err := <-execResult.Finished:
		logger.Trace("command execution finished")
		return execResult.Stdout.(io.ReadCloser), execResult.Stderr.(io.ReadCloser), err
	case <-execResult.Ctx.Done():
		logger.Error("command execution timed out", "error", execResult.Ctx.Err())
		return nil, nil, execResult.Ctx.Err()
	}
}

// executeAsync is the base implementation of the executeAsync function which executes a command asynchronously.
func (e *BaseExecutor) executeAsync(command string, stdin io.ReadCloser, timeout time.Duration, script bool) (*ExecutionResult, error) {
	var exe *exec.Cmd
	var ctx context.Context
	var cancel context.CancelFunc
	var err error
	if script {
		exe, ctx, cancel, err = e.prepareScript(command, stdin, timeout)
		if err != nil {
			return nil, err
		}
	} else {
		exe, ctx, cancel, err = e.prepareCommand(command, stdin, timeout)
		if err != nil {
			return nil, err
		}
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

	// In order to ensure that the command execution artifacts are cleaned up we need to know when the command exits.
	// We can not use exe.Wait() to do this as that will trigger a cleanup of the resources and a premature close of
	// the pipes used for stdout and stderr. Instead, we need to wait for the pipes to be closed but since we pass those
	// through to the caller we need a way to have visibility to that. We end up using a custom ReadWriteCloser here to
	// allow visibility into when the pipes are closed. This is a bit of a hack, but it is needed here instead of
	// bytes.Buffer because bytes.Buffer will return EOF if read too early before there is input to read.
	outReadWriter := internal.NewExecReadWriter(stdoutPipe)
	errReadWriter := internal.NewExecReadWriter(stderrPipe)

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
		errReadWriter.Wait()
		logger.Trace("the errReadWriter has finished")
		outReadWriter.Wait()
		logger.Trace("the outReadWriter has finished")
		exitErr := exe.Wait()
		finished <- exitErr
		logger.Trace("command finished executing", "exit", exitErr)
		if cancel != nil {
			cancel()
		}
	}()

	logger.Trace("returning ExecutionResults object")
	return &ExecutionResult{
		Stdout:   outReadWriter,
		Stderr:   errReadWriter,
		Finished: finished,
		Ctx:      ctx,
	}, nil
}

func (e *BaseExecutor) writeTempScript(scriptType ScriptType, script string) (string, error) {
	var pattern string
	switch scriptType {
	case ScriptTypePowerShell:
		pattern = "go-execute-*.ps1"
	case ScriptTypeBash:
		pattern = "go-execute-*.sh"
	case ScriptTypePython:
		pattern = "go-execute-*.py"
	}
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary script file: %w", err)
	}

	if _, err := tmpFile.Write([]byte(script)); err != nil {
		return "", fmt.Errorf("failed to write to temporary script file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary script file: %w", err)
	}

	return tmpFile.Name(), nil
}

func (e *BaseExecutor) buildScriptCommand(scriptType ScriptType, scriptPath string, arguments []string, parameters map[string]string) (command string, err error) {
	switch scriptType {
	case ScriptTypePowerShell:
		command = fmt.Sprintf("powershell.exe -NoProfile -NonInteractive -ExecutionPolicy Bypass -File %s", scriptPath)
		for key, value := range parameters {
			command += fmt.Sprintf(" -%s '%s'", key, value)
		}
	case ScriptTypeBash:
		return "", errors.New("bash support not yet implemented")
	case ScriptTypePython:
		return "", errors.New("python support not yet implemented")
	default:
		return "", errors.New("unsupported script type")
	}

	return command, nil
}

func (e *BaseExecutor) prepareScript(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout != 0 {
		logger.Trace("configuring command timeout", "timeout", timeout)
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command)
	if err != nil {
		logger.Error("failed to get command parts", "error", err)
		return nil, ctx, cancel, err
	}
	logger.Trace("split command into the parts", "cmdParts", cmdParts)

	if len(cmdParts) == 0 {
		err = errors.New("empty command")
		logger.Error("failed to get command parts", "error", err)
		return nil, ctx, cancel, err
	}

	var binary string
	var args []string
	binary, err = exec.LookPath(cmdParts[0])
	if err != nil {
		logger.Error("failed to find binary path", "error", err)
		return nil, ctx, cancel, err
	}
	logger.Trace("binary found", "binary", binary)
	args = cmdParts[1:]
	logger.Trace("setting commandcontext", "binary", binary, "args", args)

	exe := exec.CommandContext(ctx, binary, args...)
	exe.Stdin = stdin
	exe.Env = e.environment
	logger.Trace("command context set", "environment", exe.Env)

	if e.user != "" {
		err := e.configureUser(ctx, cancel, exe)
		if err != nil {
			logger.Error("failed to configure command user", "error", err)
			return exe, ctx, cancel, err
		}
		logger.Trace("configured user for execution", "user", e.User)
	}

	return exe, ctx, cancel, nil
}

// prepareCommand is the base implementation of the prepareCommand function which prepares the command for execution.
func (e *BaseExecutor) prepareCommand(command string, stdin io.ReadCloser, timeout time.Duration) (*exec.Cmd, context.Context, context.CancelFunc, error) {
	ctx := context.Background()
	var cancel context.CancelFunc

	if timeout != 0 {
		logger.Trace("configuring command timeout", "timeout", timeout)
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}

	cmdParts, err := utilities.Fields(command)
	if err != nil {
		logger.Error("failed to get command parts", "error", err)
		return nil, ctx, cancel, err
	}
	logger.Trace("split command into the parts", "cmdParts", cmdParts)

	if len(cmdParts) == 0 {
		err = errors.New("empty command")
		logger.Error("failed to get command parts", "error", err)
		return nil, ctx, cancel, err
	}

	var binary string
	var args []string
	if !e.UsingShell() {
		binary, err = exec.LookPath(cmdParts[0])
		if err != nil {
			logger.Error("failed to find binary path", "error", err)
			return nil, ctx, cancel, err
		}
		logger.Trace("binary found", "binary", binary)
		args = cmdParts[1:]
	} else {
		binary = e.shell
		switch strings.ToLower(binary) {
		case "cmd", "cmd.exe":
			args = []string{"/c", command}
		case "powershell", "powershell.exe":
			args = []string{"-NoProfile", "-NonInteractive", "-Command", command}
		default:
			args = []string{"-c", command}
		}
	}

	exe := exec.CommandContext(ctx, binary, args...)
	exe.Stdin = stdin
	exe.Env = e.environment
	logger.Trace("command context set", "environment", exe.Env)

	if e.user != "" {
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
func (e *BaseExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	return errors.New("this method must be implemented by the platform-specific executor")
}

// Struct and methods to allow basic execution without needing to instantiate a new Executor.
var defaultExecutor = NewExecutor(
	WithEnvironment(os.Environ()),
	WithDefaultShell(),
)

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
