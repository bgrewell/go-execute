package execute

import (
	"testing"
)

func TestFindBuilder(t *testing.T) {
	t.Run("creates builder with path", func(t *testing.T) {
		f := Find("/tmp")
		if f == nil {
			t.Fatal("expected non-nil builder")
		}
		if f.path != "/tmp" {
			t.Errorf("expected path '/tmp', got %q", f.path)
		}
	})

	t.Run("creates builder with no path", func(t *testing.T) {
		f := Find()
		if f.path != "" {
			t.Errorf("expected empty path, got %q", f.path)
		}
	})

	t.Run("uses default executor", func(t *testing.T) {
		f := Find()
		if f.executor == nil {
			t.Error("expected non-nil default executor")
		}
	})

	t.Run("sets path via InPath", func(t *testing.T) {
		f := Find().InPath("/home")
		if f.path != "/home" {
			t.Errorf("expected path '/home', got %q", f.path)
		}
	})

	t.Run("sets name pattern", func(t *testing.T) {
		f := Find().Name("*.txt")
		if f.name != "*.txt" {
			t.Errorf("expected name '*.txt', got %q", f.name)
		}
	})

	t.Run("sets files only", func(t *testing.T) {
		f := Find().FilesOnly()
		if f.fileType != "f" {
			t.Errorf("expected fileType 'f', got %q", f.fileType)
		}
	})

	t.Run("sets dirs only", func(t *testing.T) {
		f := Find().DirsOnly()
		if f.fileType != "d" {
			t.Errorf("expected fileType 'd', got %q", f.fileType)
		}
	})

	t.Run("sets max depth", func(t *testing.T) {
		f := Find().MaxDepth(3)
		if f.maxDepth != 3 {
			t.Errorf("expected maxDepth 3, got %d", f.maxDepth)
		}
	})

	t.Run("sets min depth", func(t *testing.T) {
		f := Find().MinDepth(1)
		if f.minDepth != 1 {
			t.Errorf("expected minDepth 1, got %d", f.minDepth)
		}
	})

	t.Run("sets newer than", func(t *testing.T) {
		f := Find().NewerThan("/tmp/ref")
		if f.newerThan != "/tmp/ref" {
			t.Errorf("expected newerThan '/tmp/ref', got %q", f.newerThan)
		}
	})

	t.Run("sets older than", func(t *testing.T) {
		f := Find().OlderThan("/tmp/ref")
		if f.olderThan != "/tmp/ref" {
			t.Errorf("expected olderThan '/tmp/ref', got %q", f.olderThan)
		}
	})

	t.Run("sets min size", func(t *testing.T) {
		f := Find().MinSize("100k")
		if f.minSize != "100k" {
			t.Errorf("expected minSize '100k', got %q", f.minSize)
		}
	})

	t.Run("sets max size", func(t *testing.T) {
		f := Find().MaxSize("1M")
		if f.maxSize != "1M" {
			t.Errorf("expected maxSize '1M', got %q", f.maxSize)
		}
	})

	t.Run("sets permissions", func(t *testing.T) {
		f := Find().Permissions("755")
		if f.permissions != "755" {
			t.Errorf("expected permissions '755', got %q", f.permissions)
		}
	})

	t.Run("accepts custom executor", func(t *testing.T) {
		mock := &MockRunner{}
		f := Find().WithExecutor(mock)
		if f.executor != mock {
			t.Error("expected custom executor to be set")
		}
	})

	t.Run("chains multiple options fluently", func(t *testing.T) {
		f := Find("/var").
			Name("*.log").
			FilesOnly().
			MaxDepth(5).
			MinSize("1k").
			MaxSize("100M")

		if f.path != "/var" {
			t.Error("path not preserved through chain")
		}
		if f.name != "*.log" {
			t.Error("name not preserved through chain")
		}
		if f.fileType != "f" {
			t.Error("fileType not preserved through chain")
		}
		if f.maxDepth != 5 {
			t.Error("maxDepth not preserved through chain")
		}
		if f.minSize != "1k" || f.maxSize != "100M" {
			t.Error("size filters not preserved through chain")
		}
	})
}

func TestFindBuilder_Run(t *testing.T) {
	t.Run("calls executor with correct command", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout:   "/tmp/file1.txt\n/tmp/file2.txt\n",
					ExitCode: 0,
				}, nil
			},
		}

		result, err := Find("/tmp").Name("*.txt").WithExecutor(mock).Run()

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

	t.Run("parses output into paths", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout: "/tmp/a.txt\n/tmp/b.txt\n/tmp/c.txt\n",
				}, nil
			},
		}

		result, err := Find("/tmp").WithExecutor(mock).Run()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Paths) != 3 {
			t.Fatalf("expected 3 paths, got %d", len(result.Paths))
		}
	})
}

func TestFindResult(t *testing.T) {
	t.Run("contains raw result", func(t *testing.T) {
		raw := &Result{Stdout: "test output"}
		result := &FindResult{
			Raw:   raw,
			Paths: nil,
		}
		if result.Raw != raw {
			t.Error("raw result not preserved")
		}
	})

	t.Run("contains parsed paths", func(t *testing.T) {
		paths := []string{"/a", "/b", "/c"}
		result := &FindResult{
			Paths: paths,
		}
		if len(result.Paths) != 3 {
			t.Errorf("expected 3 paths, got %d", len(result.Paths))
		}
	})
}
