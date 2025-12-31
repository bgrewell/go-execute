// Package parser provides utilities for parsing command strings.
package parser

import (
	"errors"
	"strings"
	"unicode"
)

// ErrUnclosedQuote is returned when a command string has mismatched quotes.
var ErrUnclosedQuote = errors.New("unclosed quote in command string")

// Fields splits a command string into fields while respecting quoted strings.
// Both single and double quotes are honored, and backslash escaping is supported.
//
// Examples:
//
//	Fields(`echo "hello world"`)     -> ["echo", `"hello world"`]
//	Fields(`grep 'pattern' file`)    -> ["grep", "'pattern'", "file"]
//	Fields(`ls -la /tmp`)            -> ["ls", "-la", "/tmp"]
//	Fields(`echo "it's fine"`)       -> ["echo", `"it's fine"`]
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
		return nil, ErrUnclosedQuote
	}

	return fields, nil
}

// StripQuotes removes surrounding quotes from a string if present.
// It handles both single and double quotes.
func StripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
