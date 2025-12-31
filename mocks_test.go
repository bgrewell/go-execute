package execute

import "context"

// MockRunner is a mock implementation of CommandRunner for testing
type MockRunner struct {
	// ExecuteFunc is called when Execute is invoked
	ExecuteFunc func(command string, args ...string) (*Result, error)

	// ExecuteContextFunc is called when ExecuteContext is invoked
	ExecuteContextFunc func(ctx context.Context, command string, args ...string) (*Result, error)

	// Calls records all calls made to the mock
	Calls []MockCall
}

// MockCall records a single call to the mock
type MockCall struct {
	Method  string
	Command string
	Args    []string
}

// Execute implements CommandRunner
func (m *MockRunner) Execute(command string, args ...string) (*Result, error) {
	m.Calls = append(m.Calls, MockCall{
		Method:  "Execute",
		Command: command,
		Args:    args,
	})

	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command, args...)
	}
	return &Result{}, nil
}

// ExecuteContext implements CommandRunner
func (m *MockRunner) ExecuteContext(ctx context.Context, command string, args ...string) (*Result, error) {
	m.Calls = append(m.Calls, MockCall{
		Method:  "ExecuteContext",
		Command: command,
		Args:    args,
	})

	if m.ExecuteContextFunc != nil {
		return m.ExecuteContextFunc(ctx, command, args...)
	}
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command, args...)
	}
	return &Result{}, nil
}

// MockExecutor is a mock implementation of the full Executor interface for testing
type MockExecutor struct {
	MockRunner

	// ExecuteAsyncFunc is called when ExecuteAsync is invoked
	ExecuteAsyncFunc func(command string, args ...string) (*AsyncResult, error)

	// ExecuteAsyncContextFunc is called when ExecuteAsyncContext is invoked
	ExecuteAsyncContextFunc func(ctx context.Context, command string, args ...string) (*AsyncResult, error)

	// Configuration storage
	env        []string
	workingDir string
	shell      string
	user       string
}

// ExecuteAsync implements Executor
func (m *MockExecutor) ExecuteAsync(command string, args ...string) (*AsyncResult, error) {
	m.Calls = append(m.Calls, MockCall{
		Method:  "ExecuteAsync",
		Command: command,
		Args:    args,
	})

	if m.ExecuteAsyncFunc != nil {
		return m.ExecuteAsyncFunc(command, args...)
	}
	return &AsyncResult{}, nil
}

// ExecuteAsyncContext implements Executor
func (m *MockExecutor) ExecuteAsyncContext(ctx context.Context, command string, args ...string) (*AsyncResult, error) {
	m.Calls = append(m.Calls, MockCall{
		Method:  "ExecuteAsyncContext",
		Command: command,
		Args:    args,
	})

	if m.ExecuteAsyncContextFunc != nil {
		return m.ExecuteAsyncContextFunc(ctx, command, args...)
	}
	if m.ExecuteAsyncFunc != nil {
		return m.ExecuteAsyncFunc(command, args...)
	}
	return &AsyncResult{}, nil
}

// SetEnvironment implements Executor
func (m *MockExecutor) SetEnvironment(env []string) Executor {
	m.env = env
	return m
}

// SetWorkingDir implements Executor
func (m *MockExecutor) SetWorkingDir(dir string) Executor {
	m.workingDir = dir
	return m
}

// SetShell implements Executor
func (m *MockExecutor) SetShell(shell string) Executor {
	m.shell = shell
	return m
}

// SetUser implements Executor
func (m *MockExecutor) SetUser(user string) Executor {
	m.user = user
	return m
}

// SetSudoPassword implements Executor
func (m *MockExecutor) SetSudoPassword(password string) Executor {
	return m
}

// Environment implements Executor
func (m *MockExecutor) Environment() []string {
	return m.env
}

// WorkingDir implements Executor
func (m *MockExecutor) WorkingDir() string {
	return m.workingDir
}

// Shell implements Executor
func (m *MockExecutor) Shell() string {
	return m.shell
}

// User implements Executor
func (m *MockExecutor) User() string {
	return m.user
}

// Compile-time interface compliance checks
var _ CommandRunner = (*MockRunner)(nil)
var _ Executor = (*MockExecutor)(nil)
