package execute

import "context"

// FindBuilder provides a fluent interface for finding files.
// It abstracts platform differences (find on Unix, dir /s on Windows).
type FindBuilder struct {
	executor CommandRunner

	path        string
	name        string // filename pattern
	fileType    string // "f" for file, "d" for directory
	maxDepth    int
	minDepth    int
	newerThan   string // file path for time comparison
	olderThan   string // file path for time comparison
	minSize     string // e.g., "+100k", "-1M"
	maxSize     string
	permissions string // e.g., "755"
}

// Find creates a new FindBuilder for searching files.
// If no path is specified, the current directory is used.
func Find(path ...string) *FindBuilder {
	b := &FindBuilder{
		executor: defaultExecutor,
	}
	if len(path) > 0 {
		b.path = path[0]
	}
	return b
}

// InPath sets the starting directory for the search
func (f *FindBuilder) InPath(path string) *FindBuilder {
	f.path = path
	return f
}

// Name sets the filename pattern to match (supports wildcards)
func (f *FindBuilder) Name(pattern string) *FindBuilder {
	f.name = pattern
	return f
}

// FilesOnly limits results to regular files
func (f *FindBuilder) FilesOnly() *FindBuilder {
	f.fileType = "f"
	return f
}

// DirsOnly limits results to directories
func (f *FindBuilder) DirsOnly() *FindBuilder {
	f.fileType = "d"
	return f
}

// MaxDepth limits the search depth
func (f *FindBuilder) MaxDepth(depth int) *FindBuilder {
	f.maxDepth = depth
	return f
}

// MinDepth sets the minimum directory depth
func (f *FindBuilder) MinDepth(depth int) *FindBuilder {
	f.minDepth = depth
	return f
}

// NewerThan finds files newer than the given reference file
func (f *FindBuilder) NewerThan(refFile string) *FindBuilder {
	f.newerThan = refFile
	return f
}

// OlderThan finds files older than the given reference file
func (f *FindBuilder) OlderThan(refFile string) *FindBuilder {
	f.olderThan = refFile
	return f
}

// MinSize finds files at least this size (e.g., "100k", "1M", "1G")
func (f *FindBuilder) MinSize(size string) *FindBuilder {
	f.minSize = size
	return f
}

// MaxSize finds files at most this size (e.g., "100k", "1M", "1G")
func (f *FindBuilder) MaxSize(size string) *FindBuilder {
	f.maxSize = size
	return f
}

// Permissions finds files with the given permission mode (e.g., "755")
func (f *FindBuilder) Permissions(mode string) *FindBuilder {
	f.permissions = mode
	return f
}

// WithExecutor sets a custom executor (useful for testing)
func (f *FindBuilder) WithExecutor(e CommandRunner) *FindBuilder {
	f.executor = e
	return f
}

// Run executes the find command and returns structured results
func (f *FindBuilder) Run() (*FindResult, error) {
	return f.RunContext(context.Background())
}

// RunContext executes the find command with context support
func (f *FindBuilder) RunContext(ctx context.Context) (*FindResult, error) {
	panic("not implemented")
}

// FindResult holds the structured output from a find command
type FindResult struct {
	// Raw contains the underlying command result
	Raw *Result

	// Paths contains the found file paths
	Paths []string
}
