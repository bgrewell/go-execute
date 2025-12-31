//go:build windows

package execute

// buildLsArgs constructs the arguments for the Windows dir command
func (l *LsBuilder) buildArgs() (command string, args []string) {
	panic("not implemented")
}

// parseLsOutput parses Windows dir output into structured entries
func parseLsOutput(output string, longFormat bool) []FileEntry {
	panic("not implemented")
}
