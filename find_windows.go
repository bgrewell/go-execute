//go:build windows

package execute

// buildFindArgs constructs the arguments for the Windows dir /s command
func (f *FindBuilder) buildArgs() (command string, args []string) {
	panic("not implemented")
}

// parseFindOutput parses Windows dir /s output into a list of paths
func parseFindOutput(output string) []string {
	panic("not implemented")
}
