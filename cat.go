package execute

import "context"

// CatBuilder provides a fluent interface for reading file contents.
// It abstracts platform differences (cat on Unix, type on Windows).
type CatBuilder struct {
	executor CommandRunner

	paths       []string
	lineNumbers bool
	showEnds    bool // show $ at end of each line
	showTabs    bool // show ^I for tabs
	squeezeBlank bool // squeeze multiple blank lines
}

// Cat creates a new CatBuilder for reading file contents
func Cat(paths ...string) *CatBuilder {
	return &CatBuilder{
		executor: defaultExecutor,
		paths:    paths,
	}
}

// Path adds a file path to read
func (c *CatBuilder) Path(path string) *CatBuilder {
	c.paths = append(c.paths, path)
	return c
}

// Paths adds multiple file paths to read
func (c *CatBuilder) Paths(paths ...string) *CatBuilder {
	c.paths = append(c.paths, paths...)
	return c
}

// WithLineNumbers includes line numbers in output
func (c *CatBuilder) WithLineNumbers() *CatBuilder {
	c.lineNumbers = true
	return c
}

// ShowEnds displays $ at the end of each line
func (c *CatBuilder) ShowEnds() *CatBuilder {
	c.showEnds = true
	return c
}

// ShowTabs displays ^I for tab characters
func (c *CatBuilder) ShowTabs() *CatBuilder {
	c.showTabs = true
	return c
}

// SqueezeBlank squeezes multiple adjacent blank lines into one
func (c *CatBuilder) SqueezeBlank() *CatBuilder {
	c.squeezeBlank = true
	return c
}

// WithExecutor sets a custom executor (useful for testing)
func (c *CatBuilder) WithExecutor(e CommandRunner) *CatBuilder {
	c.executor = e
	return c
}

// Run executes the cat command and returns structured results
func (c *CatBuilder) Run() (*CatResult, error) {
	return c.RunContext(context.Background())
}

// RunContext executes the cat command with context support
func (c *CatBuilder) RunContext(ctx context.Context) (*CatResult, error) {
	panic("not implemented")
}

// CatResult holds the output from a cat command
type CatResult struct {
	// Raw contains the underlying command result
	Raw *Result

	// Content is the concatenated file contents
	Content string
}
