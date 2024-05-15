package utilities // fields is a special function similar in nature to the standard libraries strings.Fields except that

import (
	"encoding/csv"
	"regexp"
	"strings"
)

// Fields parses a command line string into individual arguments.
// It handles Windows paths with spaces by detecting the executable path first.
func Fields(s string) ([]string, error) {
	// Regex to match Windows executable paths (modify as necessary for other extensions)
	re := regexp.MustCompile(`^([a-zA-Z]:\\[^\n\r]*?\.(exe|bat|ps1))`)

	// Find the first match for Windows executable paths
	if matches := re.FindStringSubmatch(s); len(matches) > 1 {
		executable := matches[1]
		remainingArgs := strings.TrimPrefix(s, executable)
		args, err := parseArguments(remainingArgs)
		if err != nil {
			return nil, err
		}
		return append([]string{executable}, args...), nil
	}

	// Default parsing for cases without a detected Windows path
	return parseArguments(s)
}

// parseArguments uses a CSV reader to parse the string using a space as a delimiter.
func parseArguments(input string) ([]string, error) {
	if input == "" {
		return []string, nil
	}
	r := csv.NewReader(strings.NewReader(input))
	r.Comma = ' '
	r.TrimLeadingSpace = true
	return r.Read()
}
