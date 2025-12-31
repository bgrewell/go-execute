// Package execute provides a flexible and testable interface for executing
// shell commands in Go. It offers two main APIs:
//
//   - A raw Executor interface for direct command execution with full control
//     over environment, working directory, user context, and async execution.
//
//   - Safe command wrappers (Grep, Ls, Cat, Find) that provide a builder pattern
//     API with structured results, abstracting away platform differences.
//
// # Basic Usage
//
// For simple command execution, use the package-level convenience functions:
//
//	result, err := execute.Run("ls", "-la")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Stdout)
//
// # Executor Usage
//
// For more control, create an Executor with options:
//
//	exec := execute.NewExecutor(
//	    execute.WithWorkingDir("/tmp"),
//	    execute.WithEnvironment([]string{"FOO=bar"}),
//	)
//	result, err := exec.Execute("echo", "$FOO")
//
// # Safe Command Wrappers
//
// For common commands, use the type-safe builders:
//
//	result, err := execute.Grep("pattern").
//	    InPath("/var/log").
//	    Recursive().
//	    IgnoreCase().
//	    Run()
//
//	for _, match := range result.Matches {
//	    fmt.Printf("%s:%d: %s\n", match.File, match.Line, match.Content)
//	}
//
// # Testability
//
// All command wrappers accept a CommandRunner interface, making them easy to test:
//
//	mock := &MockRunner{...}
//	result, err := execute.Grep("pattern").WithExecutor(mock).Run()
package execute

import "context"

// defaultExecutor is the package-level executor used by convenience functions
var defaultExecutor = NewExecutor()

// Run executes a command using the default executor and returns the result
func Run(command string, args ...string) (*Result, error) {
	return defaultExecutor.Execute(command, args...)
}

// RunContext executes a command with context using the default executor
func RunContext(ctx context.Context, command string, args ...string) (*Result, error) {
	return defaultExecutor.ExecuteContext(ctx, command, args...)
}

// RunAsync executes a command asynchronously using the default executor
func RunAsync(command string, args ...string) (*AsyncResult, error) {
	return defaultExecutor.ExecuteAsync(command, args...)
}

// RunAsyncContext executes a command asynchronously with context using the default executor
func RunAsyncContext(ctx context.Context, command string, args ...string) (*AsyncResult, error) {
	return defaultExecutor.ExecuteAsyncContext(ctx, command, args...)
}

// SetDefaultExecutor replaces the package-level default executor.
// This is useful for testing or for setting global configuration.
func SetDefaultExecutor(e *BaseExecutor) {
	defaultExecutor = e
}

// DefaultExecutor returns the package-level default executor
func DefaultExecutor() *BaseExecutor {
	return defaultExecutor
}
