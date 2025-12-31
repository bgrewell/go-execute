package execute

import "errors"

// Sentinel errors for common command execution failures
var (
	// ErrCommandNotFound indicates the requested command was not found in PATH
	ErrCommandNotFound = errors.New("command not found")

	// ErrTimeout indicates the command exceeded its execution timeout
	ErrTimeout = errors.New("command timed out")

	// ErrPermission indicates insufficient permissions to execute the command
	ErrPermission = errors.New("permission denied")

	// ErrKilled indicates the process was killed by a signal
	ErrKilled = errors.New("process killed")
)

// ExitError represents a command that exited with a non-zero status
type ExitError struct {
	Code   int
	Stderr string
}

// Error implements the error interface
func (e *ExitError) Error() string {
	panic("not implemented")
}
