package utilities

import (
	"errors"
	"strings"
	"unicode"
)

//// Fields is a special function similar in nature to the standard libraries strings.Fields except that
//// it honors quoted strings. It will not split inside quoted values in a string.
//func Fields(s string) (fields []string, err error) {
//	r := csv.NewReader(strings.NewReader(s))
//	r.Comma = ' '
//	return r.Read()
//}

// Fields splits a string into fields while respecting quoted strings (both single and double quotes).
func Fields(s string) ([]string, error) {
	fields := make([]string, 0)
	var field strings.Builder
	var inSingleQuote, inDoubleQuote, escaping bool

	for _, r := range s {
		switch {
		case r == '\\' && !escaping:
			escaping = true
			field.WriteRune(r)
		case r == '\\' && escaping:
			escaping = false
			field.WriteRune(r)
		case r == '\'' && !inDoubleQuote && !escaping:
			inSingleQuote = !inSingleQuote
			field.WriteRune(r)
		case r == '"' && !inSingleQuote && !escaping:
			inDoubleQuote = !inDoubleQuote
			field.WriteRune(r)
		case unicode.IsSpace(r) && !inSingleQuote && !inDoubleQuote:
			if field.Len() > 0 {
				fields = append(fields, field.String())
				field.Reset()
			}
		default:
			escaping = false
			field.WriteRune(r)
		}
	}

	// Append the last field if present
	if field.Len() > 0 {
		fields = append(fields, field.String())
	}

	if inSingleQuote || inDoubleQuote {
		return nil, errors.New("unclosed quote")
	}

	return fields, nil
}
