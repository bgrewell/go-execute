package execute

import "context"

// GrepBuilder provides a fluent interface for building grep commands.
// It abstracts platform differences (grep on Unix, findstr on Windows).
type GrepBuilder struct {
	executor CommandRunner

	pattern     string
	paths       []string
	ignoreCase  bool
	recursive   bool
	lineNumbers bool
	invertMatch bool
	wholeWord   bool
	fixedString bool
	maxCount    int
	context     int
	beforeCtx   int
	afterCtx    int
}

// Grep creates a new GrepBuilder for searching the given pattern
func Grep(pattern string) *GrepBuilder {
	return &GrepBuilder{
		executor: defaultExecutor,
		pattern:  pattern,
	}
}

// InPath adds a path to search in
func (g *GrepBuilder) InPath(path string) *GrepBuilder {
	g.paths = append(g.paths, path)
	return g
}

// InPaths adds multiple paths to search in
func (g *GrepBuilder) InPaths(paths ...string) *GrepBuilder {
	g.paths = append(g.paths, paths...)
	return g
}

// IgnoreCase enables case-insensitive matching
func (g *GrepBuilder) IgnoreCase() *GrepBuilder {
	g.ignoreCase = true
	return g
}

// Recursive enables recursive directory search
func (g *GrepBuilder) Recursive() *GrepBuilder {
	g.recursive = true
	return g
}

// WithLineNumbers includes line numbers in output
func (g *GrepBuilder) WithLineNumbers() *GrepBuilder {
	g.lineNumbers = true
	return g
}

// InvertMatch selects non-matching lines
func (g *GrepBuilder) InvertMatch() *GrepBuilder {
	g.invertMatch = true
	return g
}

// WholeWord matches only whole words
func (g *GrepBuilder) WholeWord() *GrepBuilder {
	g.wholeWord = true
	return g
}

// FixedString treats pattern as a literal string, not a regex
func (g *GrepBuilder) FixedString() *GrepBuilder {
	g.fixedString = true
	return g
}

// MaxCount stops after the specified number of matches
func (g *GrepBuilder) MaxCount(n int) *GrepBuilder {
	g.maxCount = n
	return g
}

// Context shows n lines of context around each match
func (g *GrepBuilder) Context(n int) *GrepBuilder {
	g.context = n
	return g
}

// BeforeContext shows n lines before each match
func (g *GrepBuilder) BeforeContext(n int) *GrepBuilder {
	g.beforeCtx = n
	return g
}

// AfterContext shows n lines after each match
func (g *GrepBuilder) AfterContext(n int) *GrepBuilder {
	g.afterCtx = n
	return g
}

// WithExecutor sets a custom executor (useful for testing)
func (g *GrepBuilder) WithExecutor(e CommandRunner) *GrepBuilder {
	g.executor = e
	return g
}

// Run executes the grep command and returns structured results
func (g *GrepBuilder) Run() (*GrepResult, error) {
	return g.RunContext(context.Background())
}

// RunContext executes the grep command with context support
func (g *GrepBuilder) RunContext(ctx context.Context) (*GrepResult, error) {
	panic("not implemented")
}

// GrepResult holds the structured output from a grep command
type GrepResult struct {
	// Raw contains the underlying command result
	Raw *Result

	// Matches contains the parsed grep matches
	Matches []GrepMatch
}

// GrepMatch represents a single match from grep
type GrepMatch struct {
	// File is the path to the file containing the match
	File string

	// Line is the line number of the match (if line numbers were requested)
	Line int

	// Content is the matching line content
	Content string
}
