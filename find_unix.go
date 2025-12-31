//go:build unix

package execute

// buildFindArgs constructs the arguments for the Unix find command
func (f *FindBuilder) buildArgs() (command string, args []string) {
	panic("not implemented")
}

// parseFindOutput parses Unix find output into a list of paths
func parseFindOutput(output string) []string {
	panic("not implemented")
}
