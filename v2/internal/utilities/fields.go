package utilities

import (
	"encoding/csv"
	"strings"
)

// Fields is a special function similar in nature to the standard libraries strings.Fields except that
// it honors quoted strings. It will not split inside quoted values in a string.
func Fields(s string) (fields []string, err error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	return r.Read()
}
