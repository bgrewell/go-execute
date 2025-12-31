package execute

import (
	"context"
	"os"
	"time"
)

// LsBuilder provides a fluent interface for building directory listing commands.
// It abstracts platform differences (ls on Unix, dir on Windows).
type LsBuilder struct {
	executor CommandRunner

	path       string
	all        bool // include hidden files
	long       bool // long format with details
	recursive  bool
	humanSize  bool // human-readable sizes
	sortByTime bool
	sortBySize bool
	reverse    bool
}

// Ls creates a new LsBuilder for listing directory contents.
// If no path is specified, the current directory is used.
func Ls(path ...string) *LsBuilder {
	b := &LsBuilder{
		executor: defaultExecutor,
	}
	if len(path) > 0 {
		b.path = path[0]
	}
	return b
}

// Path sets the directory path to list
func (l *LsBuilder) Path(path string) *LsBuilder {
	l.path = path
	return l
}

// All includes hidden files (those starting with .)
func (l *LsBuilder) All() *LsBuilder {
	l.all = true
	return l
}

// Long enables long format output with file details
func (l *LsBuilder) Long() *LsBuilder {
	l.long = true
	return l
}

// Recursive lists directories recursively
func (l *LsBuilder) Recursive() *LsBuilder {
	l.recursive = true
	return l
}

// HumanReadable displays file sizes in human-readable format (K, M, G)
func (l *LsBuilder) HumanReadable() *LsBuilder {
	l.humanSize = true
	return l
}

// SortByTime sorts by modification time, newest first
func (l *LsBuilder) SortByTime() *LsBuilder {
	l.sortByTime = true
	return l
}

// SortBySize sorts by file size, largest first
func (l *LsBuilder) SortBySize() *LsBuilder {
	l.sortBySize = true
	return l
}

// Reverse reverses the sort order
func (l *LsBuilder) Reverse() *LsBuilder {
	l.reverse = true
	return l
}

// WithExecutor sets a custom executor (useful for testing)
func (l *LsBuilder) WithExecutor(e CommandRunner) *LsBuilder {
	l.executor = e
	return l
}

// Run executes the ls command and returns structured results
func (l *LsBuilder) Run() (*LsResult, error) {
	return l.RunContext(context.Background())
}

// RunContext executes the ls command with context support
func (l *LsBuilder) RunContext(ctx context.Context) (*LsResult, error) {
	panic("not implemented")
}

// LsResult holds the structured output from an ls command
type LsResult struct {
	// Raw contains the underlying command result
	Raw *Result

	// Entries contains the parsed directory entries
	Entries []FileEntry
}

// FileEntry represents a single file or directory entry
type FileEntry struct {
	// Name is the base name of the file
	Name string

	// Path is the full path to the file
	Path string

	// Size is the file size in bytes
	Size int64

	// Mode is the file permission mode
	Mode os.FileMode

	// ModTime is the last modification time
	ModTime time.Time

	// IsDir indicates if this is a directory
	IsDir bool
}
