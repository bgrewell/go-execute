package execute

import (
	"testing"
)

func TestGrepBuilder(t *testing.T) {
	t.Run("creates builder with pattern", func(t *testing.T) {
		g := Grep("pattern")
		if g == nil {
			t.Fatal("expected non-nil builder")
		}
		if g.pattern != "pattern" {
			t.Errorf("expected pattern 'pattern', got %q", g.pattern)
		}
	})

	t.Run("uses default executor", func(t *testing.T) {
		g := Grep("pattern")
		if g.executor == nil {
			t.Error("expected non-nil default executor")
		}
	})

	t.Run("chains InPath calls", func(t *testing.T) {
		g := Grep("pattern").InPath("/path1").InPath("/path2")
		if len(g.paths) != 2 {
			t.Fatalf("expected 2 paths, got %d", len(g.paths))
		}
		if g.paths[0] != "/path1" || g.paths[1] != "/path2" {
			t.Errorf("unexpected paths: %v", g.paths)
		}
	})

	t.Run("InPaths adds multiple paths", func(t *testing.T) {
		g := Grep("pattern").InPaths("/a", "/b", "/c")
		if len(g.paths) != 3 {
			t.Fatalf("expected 3 paths, got %d", len(g.paths))
		}
	})

	t.Run("sets ignore case flag", func(t *testing.T) {
		g := Grep("pattern").IgnoreCase()
		if !g.ignoreCase {
			t.Error("expected ignoreCase to be true")
		}
	})

	t.Run("sets recursive flag", func(t *testing.T) {
		g := Grep("pattern").Recursive()
		if !g.recursive {
			t.Error("expected recursive to be true")
		}
	})

	t.Run("sets line numbers flag", func(t *testing.T) {
		g := Grep("pattern").WithLineNumbers()
		if !g.lineNumbers {
			t.Error("expected lineNumbers to be true")
		}
	})

	t.Run("sets invert match flag", func(t *testing.T) {
		g := Grep("pattern").InvertMatch()
		if !g.invertMatch {
			t.Error("expected invertMatch to be true")
		}
	})

	t.Run("sets whole word flag", func(t *testing.T) {
		g := Grep("pattern").WholeWord()
		if !g.wholeWord {
			t.Error("expected wholeWord to be true")
		}
	})

	t.Run("sets fixed string flag", func(t *testing.T) {
		g := Grep("pattern").FixedString()
		if !g.fixedString {
			t.Error("expected fixedString to be true")
		}
	})

	t.Run("sets max count", func(t *testing.T) {
		g := Grep("pattern").MaxCount(10)
		if g.maxCount != 10 {
			t.Errorf("expected maxCount 10, got %d", g.maxCount)
		}
	})

	t.Run("sets context lines", func(t *testing.T) {
		g := Grep("pattern").Context(3)
		if g.context != 3 {
			t.Errorf("expected context 3, got %d", g.context)
		}
	})

	t.Run("sets before context", func(t *testing.T) {
		g := Grep("pattern").BeforeContext(2)
		if g.beforeCtx != 2 {
			t.Errorf("expected beforeCtx 2, got %d", g.beforeCtx)
		}
	})

	t.Run("sets after context", func(t *testing.T) {
		g := Grep("pattern").AfterContext(2)
		if g.afterCtx != 2 {
			t.Errorf("expected afterCtx 2, got %d", g.afterCtx)
		}
	})

	t.Run("accepts custom executor", func(t *testing.T) {
		mock := &MockRunner{}
		g := Grep("pattern").WithExecutor(mock)
		if g.executor != mock {
			t.Error("expected custom executor to be set")
		}
	})

	t.Run("chains multiple options fluently", func(t *testing.T) {
		g := Grep("error").
			InPath("/var/log").
			Recursive().
			IgnoreCase().
			WithLineNumbers().
			MaxCount(100)

		if g.pattern != "error" {
			t.Error("pattern not preserved through chain")
		}
		if len(g.paths) != 1 || g.paths[0] != "/var/log" {
			t.Error("path not preserved through chain")
		}
		if !g.recursive || !g.ignoreCase || !g.lineNumbers {
			t.Error("flags not preserved through chain")
		}
		if g.maxCount != 100 {
			t.Error("maxCount not preserved through chain")
		}
	})
}

func TestGrepBuilder_Run(t *testing.T) {
	t.Run("calls executor with correct command", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				// Verify grep command structure
				return &Result{
					Stdout:   "file.txt:10:matching line\n",
					ExitCode: 0,
				}, nil
			},
		}

		result, err := Grep("pattern").
			InPath("/tmp").
			WithLineNumbers().
			WithExecutor(mock).
			Run()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if len(mock.Calls) != 1 {
			t.Errorf("expected 1 call, got %d", len(mock.Calls))
		}
	})

	t.Run("parses output into matches", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout: "file1.txt:10:first match\nfile2.txt:20:second match\n",
				}, nil
			},
		}

		result, err := Grep("pattern").
			WithLineNumbers().
			WithExecutor(mock).
			Run()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Matches) != 2 {
			t.Fatalf("expected 2 matches, got %d", len(result.Matches))
		}
		if result.Matches[0].File != "file1.txt" {
			t.Errorf("expected file 'file1.txt', got %q", result.Matches[0].File)
		}
		if result.Matches[0].Line != 10 {
			t.Errorf("expected line 10, got %d", result.Matches[0].Line)
		}
	})
}

func TestGrepResult(t *testing.T) {
	t.Run("contains raw result", func(t *testing.T) {
		raw := &Result{Stdout: "test output"}
		result := &GrepResult{
			Raw:     raw,
			Matches: nil,
		}
		if result.Raw != raw {
			t.Error("raw result not preserved")
		}
	})

	t.Run("contains parsed matches", func(t *testing.T) {
		matches := []GrepMatch{
			{File: "test.txt", Line: 1, Content: "match"},
		}
		result := &GrepResult{
			Matches: matches,
		}
		if len(result.Matches) != 1 {
			t.Errorf("expected 1 match, got %d", len(result.Matches))
		}
	})
}
