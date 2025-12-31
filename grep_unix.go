//go:build unix

package execute

// buildGrepArgs constructs the arguments for the Unix grep command
func (g *GrepBuilder) buildArgs() (command string, args []string) {
	panic("not implemented")
}

// parseGrepOutput parses Unix grep output into structured matches
func parseGrepOutput(output string, hasLineNumbers bool) []GrepMatch {
	panic("not implemented")
}
