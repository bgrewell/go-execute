//go:build windows

package execute

// buildGrepArgs constructs the arguments for the Windows findstr command
func (g *GrepBuilder) buildArgs() (command string, args []string) {
	panic("not implemented")
}

// parseGrepOutput parses Windows findstr output into structured matches
func parseGrepOutput(output string, hasLineNumbers bool) []GrepMatch {
	panic("not implemented")
}
