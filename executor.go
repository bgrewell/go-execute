package execute

import "context"

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
	panic("not implemented")
}

// ExecuteAsync starts a command asynchronously
func (e *BaseExecutor) ExecuteAsync(command string, args ...string) (*AsyncResult, error) {
	return e.ExecuteAsyncContext(context.Background(), command, args...)
}

// ExecuteAsyncContext starts a command asynchronously with context support
func (e *BaseExecutor) ExecuteAsyncContext(ctx context.Context, command string, args ...string) (*AsyncResult, error) {
	panic("not implemented")
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
