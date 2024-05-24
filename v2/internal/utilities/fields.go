package utilities // fields is a special function similar in nature to the standard libraries strings.Fields except that
import (
	"encoding/csv"
	"strings"
)

// it honors quoted strings so it will not split inside of quoted values in a string.
func Fields(s string) (fields []string, err error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = ' '
	return r.Read()
}
