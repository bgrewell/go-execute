package execute

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	t.Run("creates executor with no options", func(t *testing.T) {
		e := NewExecutor()
		if e == nil {
			t.Fatal("expected non-nil executor")
		}
	})

	t.Run("applies environment option", func(t *testing.T) {
		env := []string{"FOO=bar", "BAZ=qux"}
		e := NewExecutor(WithEnvironment(env))

		got := e.Environment()
		if len(got) != len(env) {
			t.Fatalf("expected %d env vars, got %d", len(env), len(got))
		}
		for i, v := range env {
			if got[i] != v {
				t.Errorf("env[%d]: expected %q, got %q", i, v, got[i])
			}
		}
	})

	t.Run("applies working directory option", func(t *testing.T) {
		dir := "/tmp/test"
		e := NewExecutor(WithWorkingDir(dir))

		if e.WorkingDir() != dir {
			t.Errorf("expected working dir %q, got %q", dir, e.WorkingDir())
		}
	})

	t.Run("applies shell option", func(t *testing.T) {
		shell := "/bin/zsh"
		e := NewExecutor(WithShell(shell))

		if e.Shell() != shell {
			t.Errorf("expected shell %q, got %q", shell, e.Shell())
		}
	})

	t.Run("applies user option", func(t *testing.T) {
		user := "testuser"
		e := NewExecutor(WithUser(user))

		if e.User() != user {
			t.Errorf("expected user %q, got %q", user, e.User())
		}
	})

	t.Run("applies multiple options", func(t *testing.T) {
		e := NewExecutor(
			WithWorkingDir("/tmp"),
			WithShell("/bin/bash"),
			WithUser("root"),
		)

		if e.WorkingDir() != "/tmp" {
			t.Errorf("expected working dir /tmp, got %q", e.WorkingDir())
		}
		if e.Shell() != "/bin/bash" {
			t.Errorf("expected shell /bin/bash, got %q", e.Shell())
		}
		if e.User() != "root" {
			t.Errorf("expected user root, got %q", e.User())
		}
	})
}

func TestBaseExecutor_SettersChaining(t *testing.T) {
	e := NewExecutor()

	// Test that setters return the executor for chaining
	result := e.SetEnvironment([]string{"A=1"}).
		SetWorkingDir("/tmp").
		SetShell("/bin/bash").
		SetUser("test").
		SetSudoPassword("secret")

	if result == nil {
		t.Fatal("chained setters returned nil")
	}

	// Verify all values were set
	if e.WorkingDir() != "/tmp" {
		t.Errorf("expected working dir /tmp, got %q", e.WorkingDir())
	}
}

func TestBaseExecutor_Execute(t *testing.T) {
	t.Run("executes simple command", func(t *testing.T) {
		e := NewExecutor()
		result, err := e.Execute("echo", "hello")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(result.Stdout) != "hello" {
			t.Errorf("expected stdout 'hello', got %q", result.Stdout)
		}
		if result.ExitCode != 0 {
			t.Errorf("expected exit code 0, got %d", result.ExitCode)
		}
	})

	t.Run("captures stderr", func(t *testing.T) {
		e := NewExecutor()
		result, err := e.Execute("sh", "-c", "echo error >&2")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(result.Stderr) != "error" {
			t.Errorf("expected stderr 'error', got %q", result.Stderr)
		}
	})

	t.Run("returns exit code on failure", func(t *testing.T) {
		e := NewExecutor()
		result, err := e.Execute("sh", "-c", "exit 42")

		// Implementation returns result even on non-zero exit
		if result == nil {
			t.Fatal("expected result even on non-zero exit")
		}
		if result.ExitCode != 42 {
			t.Errorf("expected exit code 42, got %d", result.ExitCode)
		}
		_ = err // Error is nil since we got a valid exit code
	})

	t.Run("records timing information", func(t *testing.T) {
		e := NewExecutor()
		before := time.Now()
		result, err := e.Execute("sleep", "0.1")
		after := time.Now()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.StartTime.Before(before) {
			t.Error("start time is before execution began")
		}
		if result.EndTime.After(after) {
			t.Error("end time is after execution returned")
		}
		if result.Duration() < 100*time.Millisecond {
			t.Errorf("expected duration >= 100ms, got %v", result.Duration())
		}
	})

	t.Run("command not found returns error", func(t *testing.T) {
		e := NewExecutor()
		_, err := e.Execute("nonexistentcommand12345")

		if err == nil {
			t.Error("expected error for non-existent command")
		}
	})

	t.Run("respects working directory", func(t *testing.T) {
		e := NewExecutor(WithWorkingDir("/tmp"))
		result, err := e.Execute("pwd")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(result.Stdout) != "/tmp" {
			t.Errorf("expected pwd to be /tmp, got %q", result.Stdout)
		}
	})

	t.Run("shell mode executes through shell", func(t *testing.T) {
		e := NewExecutor(WithShell("/bin/bash"))
		result, err := e.Execute("echo", "hello world")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(result.Stdout) != "hello world" {
			t.Errorf("expected 'hello world', got %q", result.Stdout)
		}
	})
}

func TestBaseExecutor_ExecuteContext(t *testing.T) {
	t.Run("respects context cancellation", func(t *testing.T) {
		e := NewExecutor()
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel immediately
		cancel()

		_, err := e.ExecuteContext(ctx, "sleep", "10")
		if err == nil {
			t.Error("expected error on cancelled context")
		}
	})

	t.Run("respects context timeout", func(t *testing.T) {
		e := NewExecutor()
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := e.ExecuteContext(ctx, "sleep", "10")
		if err == nil {
			t.Error("expected error on timeout")
		}
	})
}

func TestBaseExecutor_ExecuteAsync(t *testing.T) {
	t.Run("returns async result", func(t *testing.T) {
		e := NewExecutor()
		asyncResult, err := e.ExecuteAsync("echo", "hello")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if asyncResult == nil {
			t.Fatal("expected non-nil async result")
		}

		result, err := asyncResult.Wait()
		if err != nil {
			t.Fatalf("unexpected error waiting: %v", err)
		}
		if strings.TrimSpace(result.Stdout) != "hello" {
			t.Errorf("expected 'hello', got %q", result.Stdout)
		}
	})

	t.Run("done channel closes on completion", func(t *testing.T) {
		e := NewExecutor()
		asyncResult, err := e.ExecuteAsync("echo", "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Wait should close the done channel
		_, _ = asyncResult.Wait()

		select {
		case <-asyncResult.Done:
			// Expected
		default:
			t.Error("done channel should be closed after Wait()")
		}
	})
}

func TestResult_Helpers(t *testing.T) {
	t.Run("Duration calculates correctly", func(t *testing.T) {
		start := time.Now()
		end := start.Add(5 * time.Second)

		r := &Result{
			StartTime: start,
			EndTime:   end,
		}

		if r.Duration() != 5*time.Second {
			t.Errorf("expected 5s duration, got %v", r.Duration())
		}
	})

	t.Run("Success returns true for exit code 0", func(t *testing.T) {
		r := &Result{ExitCode: 0}
		if !r.Success() {
			t.Error("expected Success() to be true for exit code 0")
		}
	})

	t.Run("Success returns false for non-zero exit code", func(t *testing.T) {
		r := &Result{ExitCode: 1}
		if r.Success() {
			t.Error("expected Success() to be false for exit code 1")
		}
	})
}

func TestDirectVsShellExecution(t *testing.T) {
	t.Run("direct mode prevents shell injection", func(t *testing.T) {
		e := NewExecutor()
		// In direct mode, this should be treated as a literal argument
		result, err := e.Execute("echo", "hello; echo injected")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// The entire argument should be printed as one line, including the semicolon
		// If shell injection worked, we'd see "injected" on a separate line
		stdout := strings.TrimSpace(result.Stdout)
		if stdout != "hello; echo injected" {
			t.Errorf("expected literal 'hello; echo injected', got %q", stdout)
		}
	})

	t.Run("shell mode allows shell features", func(t *testing.T) {
		e := NewExecutor(WithShell("/bin/bash"))
		result, err := e.Execute("echo hello && echo world")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// In shell mode, && should work
		if !strings.Contains(result.Stdout, "hello") || !strings.Contains(result.Stdout, "world") {
			t.Errorf("expected both hello and world in output, got %q", result.Stdout)
		}
	})
}
