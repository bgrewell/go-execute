package execute

import (
	"testing"
)

func TestLsBuilder(t *testing.T) {
	t.Run("creates builder with no path", func(t *testing.T) {
		l := Ls()
		if l == nil {
			t.Fatal("expected non-nil builder")
		}
		if l.path != "" {
			t.Errorf("expected empty path, got %q", l.path)
		}
	})

	t.Run("creates builder with path", func(t *testing.T) {
		l := Ls("/tmp")
		if l.path != "/tmp" {
			t.Errorf("expected path '/tmp', got %q", l.path)
		}
	})

	t.Run("uses default executor", func(t *testing.T) {
		l := Ls()
		if l.executor == nil {
			t.Error("expected non-nil default executor")
		}
	})

	t.Run("sets path via method", func(t *testing.T) {
		l := Ls().Path("/home")
		if l.path != "/home" {
			t.Errorf("expected path '/home', got %q", l.path)
		}
	})

	t.Run("sets all flag", func(t *testing.T) {
		l := Ls().All()
		if !l.all {
			t.Error("expected all to be true")
		}
	})

	t.Run("sets long format flag", func(t *testing.T) {
		l := Ls().Long()
		if !l.long {
			t.Error("expected long to be true")
		}
	})

	t.Run("sets recursive flag", func(t *testing.T) {
		l := Ls().Recursive()
		if !l.recursive {
			t.Error("expected recursive to be true")
		}
	})

	t.Run("sets human readable flag", func(t *testing.T) {
		l := Ls().HumanReadable()
		if !l.humanSize {
			t.Error("expected humanSize to be true")
		}
	})

	t.Run("sets sort by time flag", func(t *testing.T) {
		l := Ls().SortByTime()
		if !l.sortByTime {
			t.Error("expected sortByTime to be true")
		}
	})

	t.Run("sets sort by size flag", func(t *testing.T) {
		l := Ls().SortBySize()
		if !l.sortBySize {
			t.Error("expected sortBySize to be true")
		}
	})

	t.Run("sets reverse flag", func(t *testing.T) {
		l := Ls().Reverse()
		if !l.reverse {
			t.Error("expected reverse to be true")
		}
	})

	t.Run("accepts custom executor", func(t *testing.T) {
		mock := &MockRunner{}
		l := Ls().WithExecutor(mock)
		if l.executor != mock {
			t.Error("expected custom executor to be set")
		}
	})

	t.Run("chains multiple options fluently", func(t *testing.T) {
		l := Ls("/var/log").
			All().
			Long().
			HumanReadable().
			SortByTime().
			Reverse()

		if l.path != "/var/log" {
			t.Error("path not preserved through chain")
		}
		if !l.all || !l.long || !l.humanSize {
			t.Error("flags not preserved through chain")
		}
		if !l.sortByTime || !l.reverse {
			t.Error("sort options not preserved through chain")
		}
	})
}

func TestLsBuilder_Run(t *testing.T) {
	t.Run("calls executor with correct command", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout:   "file1.txt\nfile2.txt\n",
					ExitCode: 0,
				}, nil
			},
		}

		result, err := Ls("/tmp").WithExecutor(mock).Run()

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

	t.Run("parses output into file entries", func(t *testing.T) {
		t.Skip("not implemented")

		mock := &MockRunner{
			ExecuteFunc: func(command string, args ...string) (*Result, error) {
				return &Result{
					Stdout: "file1.txt\nfile2.txt\ndir1\n",
				}, nil
			},
		}

		result, err := Ls().WithExecutor(mock).Run()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result.Entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(result.Entries))
		}
	})
}

func TestLsResult(t *testing.T) {
	t.Run("contains raw result", func(t *testing.T) {
		raw := &Result{Stdout: "test output"}
		result := &LsResult{
			Raw:     raw,
			Entries: nil,
		}
		if result.Raw != raw {
			t.Error("raw result not preserved")
		}
	})

	t.Run("contains parsed entries", func(t *testing.T) {
		entries := []FileEntry{
			{Name: "test.txt", Size: 100, IsDir: false},
		}
		result := &LsResult{
			Entries: entries,
		}
		if len(result.Entries) != 1 {
			t.Errorf("expected 1 entry, got %d", len(result.Entries))
		}
	})
}

func TestFileEntry(t *testing.T) {
	t.Run("stores file metadata", func(t *testing.T) {
		entry := FileEntry{
			Name:  "test.txt",
			Path:  "/tmp/test.txt",
			Size:  1024,
			IsDir: false,
		}

		if entry.Name != "test.txt" {
			t.Errorf("expected name 'test.txt', got %q", entry.Name)
		}
		if entry.Size != 1024 {
			t.Errorf("expected size 1024, got %d", entry.Size)
		}
		if entry.IsDir {
			t.Error("expected IsDir to be false")
		}
	})
}
