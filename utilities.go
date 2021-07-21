package execute

import (
	"encoding/csv"
	"strings"
)

// Fields is a special function similar in nature to the standard libraries strings.Fields except that
// it honors quoted strings so it will not split inside of quoted values in a string.
func Fields(s string) (fields []string, err error) {
	s = strings.ReplaceAll(s,"\"", "\\\"")
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	return r.Read()
}
