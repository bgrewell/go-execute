package execute

import (
	"testing"
)

func TestCatBuilder(t *testing.T) {
	t.Run("creates builder with paths", func(t *testing.T) {
		c := Cat("/tmp/file1.txt", "/tmp/file2.txt")
		if c == nil {
			t.Fatal("expected non-nil builder")
		}
		if len(c.paths) != 2 {
			t.Fatalf("expected 2 paths, got %d", len(c.paths))
		}
	})

	t.Run("creates builder with no paths", func(t *testing.T) {
		c := Cat()
		if len(c.paths) != 0 {
			t.Errorf("expected 0 paths, got %d", len(c.paths))
		}
	})

	t.Run("uses default executor", func(t *testing.T) {
		c := Cat()
		if c.executor == nil {
			t.Error("expected non-nil default executor")
		}
	})

	t.Run("adds path via method", func(t *testing.T) {
		c := Cat().Path("/tmp/file.txt")
		if len(c.paths) != 1 {
			t.Fatalf("expected 1 path, got %d", len(c.paths))
		}
		if c.paths[0] != "/tmp/file.txt" {
			t.Errorf("expected path '/tmp/file.txt', got %q", c.paths[0])
		}
	})

	t.Run("adds multiple paths via method", func(t *testing.T) {
		c := Cat().Paths("/a", "/b", "/c")
		if len(c.paths) != 3 {
			t.Fatalf("expected 3 paths, got %d", len(c.paths))
		}
	})

	t.Run("sets line numbers flag", func(t *testing.T) {
		c := Cat().WithLineNumbers()
		if !c.lineNumbers {
			t.Error("expected lineNumbers to be true")
		}
	})

	t.Run("sets show ends flag", func(t *testing.T) {
		c := Cat().ShowEnds()
		if !c.showEnds {
			t.Error("expected showEnds to be true")
		}
	})

	t.Run("sets show tabs flag", func(t *testing.T) {
		c := Cat().ShowTabs()
		if !c.showTabs {
			t.Error("expected showTabs to be true")
		}
	})

	t.Run("sets squeeze blank flag", func(t *testing.T) {
		c := Cat().SqueezeBlank()
		if !c.squeezeBlank {
			t.Error("expected squeezeBlank to be true")
		}
	})

	t.Run("accepts custom executor", func(t *testing.T) {
		mock := &MockRunner{}
		c := Cat().WithExecutor(mock)
		if c.executor != mock {
			t.Error("expected custom executor to be set")
		}
	})

	t.Run("chains multiple options fluently", func(t *testing.T) {
		c := Cat("/etc/passwd").
			Path("/etc/group").
			WithLineNumbers().
			ShowEnds()

		if len(c.paths) != 2 {
			t.Errorf("expected 2 paths, got %d", len(c.paths))
		}
		if !c.lineNumbers || !c.showEnds {
			t.Error("flags not preserved through chain")
		}
	})
}

func TestCatBuilder_Run(t *testing.T) {
	t.Run("calls executor with correct command", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout:   "file contents here\n",
					ExitCode: 0,
				}, nil
			},
		}

		result, err := Cat("/tmp/test.txt").WithExecutor(mock).Run()

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

	t.Run("returns file content", func(t *testing.T) {
		t.Skip("not implemented")

		expectedContent := "line 1\nline 2\nline 3\n"
		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout: expectedContent,
				}, nil
			},
		}

		result, err := Cat("/tmp/test.txt").WithExecutor(mock).Run()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Content != expectedContent {
			t.Errorf("expected content %q, got %q", expectedContent, result.Content)
		}
	})
}

func TestCatResult(t *testing.T) {
	t.Run("contains raw result", func(t *testing.T) {
		raw := &Result{Stdout: "test output"}
		result := &CatResult{
			Raw:     raw,
			Content: "test output",
		}
		if result.Raw != raw {
			t.Error("raw result not preserved")
		}
	})

	t.Run("contains content", func(t *testing.T) {
		result := &CatResult{
			Content: "hello world",
		}
		if result.Content != "hello world" {
			t.Errorf("expected content 'hello world', got %q", result.Content)
		}
	})
}
