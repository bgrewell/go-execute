package execute

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bgrewell/go-execute/internal/parser"
	"github.com/bgrewell/go-execute/internal/pipe"
)

// Executor defines the full contract for command execution.
// It provides both synchronous and asynchronous execution methods,
// along with configuration for environment, user context, and shell.
type Executor interface {
	CommandRunner

	// ExecuteAsync starts a command asynchronously and returns handles to its output
	ExecuteAsync(command string, args ...string) (*AsyncResult, error)

	// ExecuteAsyncContext starts a command asynchronously with context support
	ExecuteAsyncContext(ctx context.Context, command string, args ...string) (*AsyncResult, error)

	// Configuration methods (return Executor for chaining)
	SetEnvironment(env []string) Executor
	SetWorkingDir(dir string) Executor
	SetShell(shell string) Executor
	SetUser(user string) Executor
	SetSudoPassword(password string) Executor

	// Getters for configuration inspection and testing
	Environment() []string
	WorkingDir() string
	Shell() string
	User() string
}

// CommandRunner is a minimal interface for command execution.
// Safe command wrappers depend on this interface rather than the full Executor,
// making them easy to test with mocks.
type CommandRunner interface {
	// Execute runs a command and waits for it to complete
	Execute(command string, args ...string) (*Result, error)

	// ExecuteContext runs a command with context support for cancellation/timeout
	ExecuteContext(ctx context.Context, command string, args ...string) (*Result, error)
}

// BaseExecutor is the default implementation of the Executor interface.
// It provides platform-agnostic configuration storage with platform-specific
// execution handled by OS-specific files.
type BaseExecutor struct {
	env          []string
	workingDir   string
	shell        string
	user         string
	sudoPassword string
}

// NewExecutor creates a new BaseExecutor with the given options
func NewExecutor(opts ...Option) *BaseExecutor {
	e := &BaseExecutor{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Execute runs a command and waits for it to complete
func (e *BaseExecutor) Execute(command string, args ...string) (*Result, error) {
	return e.ExecuteContext(context.Background(), command, args...)
}

// ExecuteContext runs a command with context support
func (e *BaseExecutor) ExecuteContext(ctx context.Context, command string, args ...string) (*Result, error) {
	asyncResult, err := e.ExecuteAsyncContext(ctx, command, args...)
	if err != nil {
		return nil, err
	}
	return asyncResult.Wait()
}

// ExecuteAsync starts a command asynchronously
func (e *BaseExecutor) ExecuteAsync(command string, args ...string) (*AsyncResult, error) {
	return e.ExecuteAsyncContext(context.Background(), command, args...)
}

// ExecuteAsyncContext starts a command asynchronously with context support
func (e *BaseExecutor) ExecuteAsyncContext(ctx context.Context, command string, args ...string) (*AsyncResult, error) {
	cmd, err := e.prepareCommand(ctx, command, args...)
	if err != nil {
		logger.Error("failed to prepare command", "error", err)
		return nil, err
	}

	// Set up stdout and stderr pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("failed to create stdout pipe", "error", err)
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("failed to create stderr pipe", "error", err)
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Wrap pipes in our Reader to handle timing issues
	stdoutReader := pipe.NewReader(stdoutPipe)
	stderrReader := pipe.NewReader(stderrPipe)

	startTime := time.Now()

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.Error("failed to start command", "error", err)
		return nil, fmt.Errorf("failed to start command: %w", err)
	}
	logger.Debug("started command", "command", command, "args", args)

	// Create done channel
	done := make(chan struct{})

	// Create the wait function
	waitFunc := func() (*Result, error) {
		// Wait for both pipes to be fully consumed
		stderrReader.Wait()
		stdoutReader.Wait()

		// Now wait for the command to exit
		exitErr := cmd.Wait()
		close(done)

		endTime := time.Now()

		result := &Result{
			Stdout:    stdoutReader.String(),
			Stderr:    stderrReader.String(),
			StartTime: startTime,
			EndTime:   endTime,
			Command:   command,
			Args:      args,
		}

		// Check if context was canceled/timed out
		if ctx.Err() != nil {
			logger.Debug("command terminated due to context", "error", ctx.Err())
			return result, ctx.Err()
		}

		if exitErr != nil {
			if exitError, ok := exitErr.(*exec.ExitError); ok {
				result.ExitCode = exitError.ExitCode()
			} else {
				return result, exitErr
			}
		}

		logger.Debug("command completed", "command", command, "exitCode", result.ExitCode)
		return result, nil
	}

	return &AsyncResult{
		Stdout: stdoutReader,
		Stderr: stderrReader,
		Done:   done,
		Wait:   waitFunc,
	}, nil
}

// prepareCommand creates an exec.Cmd based on the executor configuration.
// If a shell is configured, the command is executed through the shell.
// Otherwise, the command is executed directly for safety.
func (e *BaseExecutor) prepareCommand(ctx context.Context, command string, args ...string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	if e.shell != "" {
		// Shell mode: execute through shell
		cmd = e.prepareShellCommand(ctx, command, args...)
	} else {
		// Direct mode (safe): resolve and execute binary directly
		var err error
		cmd, err = e.prepareDirectCommand(ctx, command, args...)
		if err != nil {
			return nil, err
		}
	}

	// Apply common configuration
	cmd.Env = e.env
	cmd.Dir = e.workingDir

	// Configure user if specified (platform-specific)
	if e.user != "" {
		if err := e.configureUser(cmd); err != nil {
			logger.Error("failed to configure user", "error", err)
			return nil, fmt.Errorf("failed to configure user: %w", err)
		}
	}

	return cmd, nil
}

// prepareDirectCommand creates an exec.Cmd for direct execution (no shell).
// This is the safe default that prevents shell injection.
func (e *BaseExecutor) prepareDirectCommand(ctx context.Context, command string, args ...string) (*exec.Cmd, error) {
	// Handle sudo with password
	if command == "sudo" && e.sudoPassword != "" {
		return e.prepareSudoCommand(ctx, args...)
	}

	// Look up the binary path
	binary, err := exec.LookPath(command)
	if err != nil {
		logger.Error("command not found", "command", command, "error", err)
		return nil, fmt.Errorf("%w: %s", ErrCommandNotFound, command)
	}
	logger.Debug("resolved binary path", "command", command, "binary", binary)

	return exec.CommandContext(ctx, binary, args...), nil
}

// prepareSudoCommand handles sudo commands with password input
func (e *BaseExecutor) prepareSudoCommand(ctx context.Context, args ...string) (*exec.Cmd, error) {
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return nil, fmt.Errorf("%w: sudo", ErrCommandNotFound)
	}

	// Prepend -S flag to read password from stdin
	sudoArgs := append([]string{"-S"}, args...)
	cmd := exec.CommandContext(ctx, sudoPath, sudoArgs...)

	// Provide password via stdin
	cmd.Stdin = strings.NewReader(e.sudoPassword + "\n")

	return cmd, nil
}

// prepareShellCommand creates an exec.Cmd for shell execution.
// This is used when WithShell() is configured.
func (e *BaseExecutor) prepareShellCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	// Build the full command string
	fullCommand := command
	if len(args) > 0 {
		// Quote args that contain spaces
		quotedArgs := make([]string, len(args))
		for i, arg := range args {
			if strings.Contains(arg, " ") && !strings.HasPrefix(arg, "'") && !strings.HasPrefix(arg, "\"") {
				quotedArgs[i] = fmt.Sprintf("%q", arg)
			} else {
				quotedArgs[i] = arg
			}
		}
		fullCommand = command + " " + strings.Join(quotedArgs, " ")
	}

	// Handle sudo with password in shell mode
	if strings.Contains(fullCommand, "sudo ") && e.sudoPassword != "" {
		// Replace sudo with echo password | sudo -S
		fullCommand = strings.Replace(fullCommand, "sudo ", fmt.Sprintf("echo '%s' | sudo -S ", e.sudoPassword), -1)
	}

	// Build shell arguments based on shell type
	var shellArgs []string
	shellLower := strings.ToLower(e.shell)

	switch {
	case strings.Contains(shellLower, "cmd"):
		shellArgs = []string{"/c", fullCommand}
	case strings.Contains(shellLower, "powershell"):
		shellArgs = []string{"-NoProfile", "-NonInteractive", "-Command", fullCommand}
	default:
		// Unix shells (bash, sh, zsh, etc.)
		shellArgs = []string{"-c", fullCommand}
	}

	logger.Debug("prepared shell command", "shell", e.shell, "command", fullCommand)
	return exec.CommandContext(ctx, e.shell, shellArgs...)
}

// UsingShell returns true if a shell is configured
func (e *BaseExecutor) UsingShell() bool {
	return e.shell != ""
}

// SetEnvironment sets the environment variables and returns the executor for chaining
func (e *BaseExecutor) SetEnvironment(env []string) Executor {
	e.env = env
	return e
}

// SetWorkingDir sets the working directory and returns the executor for chaining
func (e *BaseExecutor) SetWorkingDir(dir string) Executor {
	e.workingDir = dir
	return e
}

// SetShell sets the shell and returns the executor for chaining
func (e *BaseExecutor) SetShell(shell string) Executor {
	e.shell = shell
	return e
}

// SetUser sets the user and returns the executor for chaining
func (e *BaseExecutor) SetUser(user string) Executor {
	e.user = user
	return e
}

// SetSudoPassword sets the sudo password and returns the executor for chaining
func (e *BaseExecutor) SetSudoPassword(password string) Executor {
	e.sudoPassword = password
	return e
}

// Environment returns the configured environment variables
func (e *BaseExecutor) Environment() []string {
	return e.env
}

// WorkingDir returns the configured working directory
func (e *BaseExecutor) WorkingDir() string {
	return e.workingDir
}

// Shell returns the configured shell
func (e *BaseExecutor) Shell() string {
	return e.shell
}

// User returns the configured user
func (e *BaseExecutor) User() string {
	return e.user
}

// Compile-time interface compliance check
var _ Executor = (*BaseExecutor)(nil)
var _ CommandRunner = (*BaseExecutor)(nil)

// Keep parser package imported for future use
var _ = parser.Fields
