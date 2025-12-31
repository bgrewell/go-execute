//go:build unix

package execute

// buildLsArgs constructs the arguments for the Unix ls command
func (l *LsBuilder) buildArgs() (command string, args []string) {
	panic("not implemented")
}

// parseLsOutput parses Unix ls output into structured entries
func parseLsOutput(output string, longFormat bool) []FileEntry {
	panic("not implemented")
}
