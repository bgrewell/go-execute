package parser

import (
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "simple command",
			input: "ls -la /tmp",
			want:  []string{"ls", "-la", "/tmp"},
		},
		{
			name:  "double quoted string",
			input: `echo "hello world"`,
			want:  []string{"echo", `"hello world"`},
		},
		{
			name:  "single quoted string",
			input: `grep 'pattern' file`,
			want:  []string{"grep", "'pattern'", "file"},
		},
		{
			name:  "mixed quotes",
			input: `echo "it's fine"`,
			want:  []string{"echo", `"it's fine"`},
		},
		{
			name:  "multiple spaces",
			input: "ls    -la    /tmp",
			want:  []string{"ls", "-la", "/tmp"},
		},
		{
			name:  "leading and trailing spaces",
			input: "  ls -la  ",
			want:  []string{"ls", "-la"},
		},
		{
			name:  "empty string",
			input: "",
			want:  []string{},
		},
		{
			name:  "only spaces",
			input: "   ",
			want:  []string{},
		},
		{
			name:  "escaped quote in double quotes",
			input: `echo "hello \"world\""`,
			want:  []string{"echo", `"hello \"world\""`},
		},
		{
			name:    "unclosed double quote",
			input:   `echo "hello`,
			wantErr: true,
		},
		{
			name:    "unclosed single quote",
			input:   `echo 'hello`,
			wantErr: true,
		},
		{
			name:  "complex command",
			input: `find /var/log -name "*.log" -type f`,
			want:  []string{"find", "/var/log", "-name", `"*.log"`, "-type", "f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Fields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStripQuotes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"hello"`, "hello"},
		{`'hello'`, "hello"},
		{`hello`, "hello"},
		{`"hello`, `"hello`},
		{`hello"`, `hello"`},
		{`""`, ""},
		{`''`, ""},
		{`"`, `"`},
		{``, ``},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := StripQuotes(tt.input)
			if got != tt.want {
				t.Errorf("StripQuotes(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
