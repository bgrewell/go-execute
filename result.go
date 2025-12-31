package execute

import (
	"io"
	"time"
)

// Result holds the output from a completed command execution
type Result struct {
	// Stdout contains the standard output from the command
	Stdout string

	// Stderr contains the standard error output from the command
	Stderr string

	// ExitCode is the exit status of the command
	ExitCode int

	// StartTime is when the command began execution
	StartTime time.Time

	// EndTime is when the command finished execution
	EndTime time.Time

	// Command is the executable that was run
	Command string

	// Args are the arguments passed to the command
	Args []string
}

// Duration returns the execution duration of the command
func (r *Result) Duration() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

// Success returns true if the command exited with code 0
func (r *Result) Success() bool {
	return r.ExitCode == 0
}

// AsyncResult holds handles to a running asynchronous command
type AsyncResult struct {
	// Stdout provides streaming access to standard output
	Stdout io.ReadCloser

	// Stderr provides streaming access to standard error
	Stderr io.ReadCloser

	// Done is closed when the command completes
	Done <-chan struct{}

	// Wait blocks until the command completes and returns the final result
	Wait func() (*Result, error)
}
