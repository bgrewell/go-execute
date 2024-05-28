package execute

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestExecuteReturnsCombinedOutput(t *testing.T) {
	command := "echo Hello, World!"
	combined, err := Execute(command)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	combined = strings.Replace(combined, "\r", "", -1)
	if combined != "Hello, World!\n" {
		t.Fatalf("Unexpected combined output: %s", combined)
	}
}

func TestExecuteHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	_, err := Execute(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteHandlesEmptyCommand(t *testing.T) {
	command := ""
	_, err := Execute(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteSeparateReturnsOutput(t *testing.T) {
	command := "echo Hello, World!"
	stdout, stderr, err := ExecuteSeparate(command)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout = strings.Replace(stdout, "\r", "", -1)
	if stdout != "Hello, World!\n" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	if stderr != "" {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
}

func TestExecuteSeparateHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	_, _, err := ExecuteSeparate(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteSeparateHandlesEmptyCommand(t *testing.T) {
	command := ""
	_, _, err := ExecuteSeparate(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteAsyncReturnsResult(t *testing.T) {
	command := "echo Hello, World!"
	execResult, err := ExecuteAsync(command)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	stdout = []byte(strings.Replace(string(stdout), "\r", "", -1))
	if string(stdout) != "Hello, World!\n" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
}

func TestExecuteAsyncHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	execResult, err := ExecuteAsync(command)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	if string(stdout) != "" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	stderr, err := io.ReadAll(execResult.Stderr)
	if err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	}
	if string(stderr) == "" {
		t.Fatalf("Expected stderr, got empty")
	}
}

func TestExecuteAsyncHandlesEmptyCommand(t *testing.T) {
	command := ""
	execResult, err := ExecuteAsync(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if execResult != nil {
		t.Fatalf("Expected nil ExecutionResult, got %v", execResult)
	}
}

func TestExecuteAsyncWithInputReturnsResult(t *testing.T) {
	command := "cat"
	stdin := io.NopCloser(strings.NewReader("Hello, World!"))
	execResult, err := ExecuteAsyncWithInput(command, stdin)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	if string(stdout) != "Hello, World!" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
}

func TestExecuteAsyncWithInputHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	stdin := io.NopCloser(strings.NewReader(""))
	execResult, err := ExecuteAsyncWithInput(command, stdin)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	if string(stdout) != "" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	stderr, err := io.ReadAll(execResult.Stderr)
	if err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	}
	if string(stderr) == "" {
		t.Fatalf("Expected stderr, got empty")
	}
}

func TestExecuteAsyncWithInputHandlesEmptyCommand(t *testing.T) {
	command := ""
	stdin := io.NopCloser(strings.NewReader(""))
	execResult, err := ExecuteAsyncWithInput(command, stdin)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if execResult != nil {
		t.Fatalf("Expected nil ExecutionResult, got %v", execResult)
	}
}

func TestExecuteAsyncWithInputHandlesNilStdin(t *testing.T) {
	command := "echo Hello, World!"
	execResult, err := ExecuteAsyncWithInput(command, nil)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	stdout = []byte(strings.Replace(string(stdout), "\r", "", -1))
	if string(stdout) != "Hello, World!\n" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
}

func TestExecuteWithTimeoutReturnsCombinedOutput(t *testing.T) {
	command := "echo Hello, World!"
	combined, err := ExecuteWithTimeout(command, 3*time.Second)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	combined = strings.Replace(combined, "\r", "", -1)
	if combined != "Hello, World!\n" {
		t.Fatalf("Unexpected combined output: %s", combined)
	}
}

func TestExecuteWithTimeoutHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	_, err := ExecuteWithTimeout(command, 3*time.Second)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteWithTimeoutHandlesEmptyCommand(t *testing.T) {
	command := ""
	_, err := ExecuteWithTimeout(command, 3*time.Second)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteWithTimeoutExceeds(t *testing.T) {
	command := "sleep 5"
	_, err := ExecuteWithTimeout(command, 3*time.Second)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteSeparateWithTimeoutReturnsOutput(t *testing.T) {
	command := "echo Hello, World!"
	stdout, stderr, err := ExecuteSeparateWithTimeout(command, 3*time.Second)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	stdout = strings.Replace(stdout, "\r", "", -1)
	if stdout != "Hello, World!\n" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	if stderr != "" {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
}

func TestExecuteSeparateWithTimeoutHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	stdout, stderr, err := ExecuteSeparateWithTimeout(command, 3*time.Second)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if stdout != "" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	if stderr != "" {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
}

func TestExecuteSeparateWithTimeoutHandlesEmptyCommand(t *testing.T) {
	command := ""
	stdout, stderr, err := ExecuteSeparateWithTimeout(command, 3*time.Second)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if stdout != "" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	if stderr != "" {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
}

func TestExecuteSeparateWithTimeoutExceeds(t *testing.T) {
	command := "sleep 5"
	stdout, stderr, err := ExecuteSeparateWithTimeout(command, 3*time.Second)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if stdout != "" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
	if stderr != "" {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
}

func TestExecuteAsyncWithTimeoutFailure(t *testing.T) {
	e := NewExecutor(
		WithShell("powershell.exe"),
	)

	command := "sleep 5"
	execResult, err := e.ExecuteAsyncWithTimeout(command, 3*time.Second)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	err = <-execResult.Finished
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	stderr, err := io.ReadAll(execResult.Stderr)
	if err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	}
	if len(stderr) > 0 {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
	if len(stdout) > 0 {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
}

func TestExecuteAsyncWithTimeoutInvalidCommand(t *testing.T) {
	command := "invalid_command"
	execResult, err := ExecuteAsyncWithTimeout(command, 3*time.Second)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	err = <-execResult.Finished
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	stderr, err := io.ReadAll(execResult.Stderr)
	if err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	}
	if len(stderr) == 0 {
		t.Fatalf("Expected stderr, got empty")
	}
	if len(stdout) > 0 {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
}

func TestExecuteAsyncWithTimeoutZeroTimeout(t *testing.T) {
	e := NewExecutor(
		WithShell("powershell.exe"),
	)

	command := "sleep 5"
	execResult, err := e.ExecuteAsyncWithTimeout(command, 0)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	err = <-execResult.Finished
	if err != nil {
		t.Fatalf("Unexpected error executing command: %v", err)
	}
	stdout, err := io.ReadAll(execResult.Stdout)
	if err != nil {
		t.Fatalf("Error reading stdout: %v", err)
	}
	stderr, err := io.ReadAll(execResult.Stderr)
	if err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	}
	if len(stderr) > 0 {
		t.Fatalf("Unexpected stderr: %s", stderr)
	}
	t.Logf("%s: %s", command, stdout)
}

func TestExecuteTTYReturnsNoErrorForValidCommand(t *testing.T) {
	command := "echo Hello, World!"
	err := ExecuteTTY(command)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
}

func TestExecuteTTYReturnsErrorForInvalidCommand(t *testing.T) {
	command := "invalid_command"
	err := ExecuteTTY(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteTTYReturnsErrorForEmptyCommand(t *testing.T) {
	command := ""
	err := ExecuteTTY(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
}

func TestExecuteScriptFromString(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("ExecuteScriptFromString_WithValidScriptAndParameters", func(t *testing.T) {
		script := "Write-Output $param1"
		parameters := map[string]string{"param1": "Hello, World!"}
		stdout, stderr, err := executor.ExecuteScriptFromString(ScriptTypePowerShell, script, nil, parameters)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if stdout != "Hello, World!" || stderr != "" {
			t.Errorf("Expected stdout 'Hello, World!', stderr '', but got stdout '%s', stderr '%s'", stdout, stderr)
		}
	})

	t.Run("ExecuteScriptFromString_WithInvalidScript", func(t *testing.T) {
		script := "Invalid-Command"
		parameters := map[string]string{}
		stdout, stderr, err := executor.ExecuteScriptFromString(ScriptTypePowerShell, script, nil, parameters)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if stdout != "" {
			t.Errorf("Expected stdout '', but got stdout '%s'", stdout)
		}

		if stderr == "" {
			t.Errorf("Expected stderr, but got empty")
		}
	})

	t.Run("ExecuteScriptFromString_WithEmptyScript", func(t *testing.T) {
		script := ""
		parameters := map[string]string{}
		stdout, stderr, err := executor.ExecuteScriptFromString(ScriptTypePowerShell, script, nil, parameters)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if stdout != "" {
			t.Errorf("Expected stdout '', but got stdout '%s'", stdout)
		}
		if stderr != "" {
			t.Errorf("Expected stderr '', but got stderr '%s'", stderr)
		}
	})
}

func TestExecuteScriptFromFile(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("ExecuteScriptFromFile_WithValidScriptAndParameters", func(t *testing.T) {
		scriptPath := "/path/to/valid/script"
		parameters := map[string]string{"param1": "Hello, World!"}
		stdout, stderr, err := executor.ExecuteScriptFromFile(ScriptTypePowerShell, scriptPath, nil, parameters)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if stdout != "Hello, World!" || stderr != "" {
			t.Errorf("Expected stdout 'Hello, World!', stderr '', but got stdout '%s', stderr '%s'", stdout, stderr)
		}
	})

	t.Run("ExecuteScriptFromFile_WithInvalidScript", func(t *testing.T) {
		scriptPath := "/path/to/invalid/script"
		parameters := map[string]string{}
		_, _, err := executor.ExecuteScriptFromFile(ScriptTypePowerShell, scriptPath, nil, parameters)

		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})

	t.Run("ExecuteScriptFromFile_WithNonexistentScript", func(t *testing.T) {
		scriptPath := "/path/to/nonexistent/script"
		parameters := map[string]string{}
		_, _, err := executor.ExecuteScriptFromFile(ScriptTypePowerShell, scriptPath, nil, parameters)

		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})
}
