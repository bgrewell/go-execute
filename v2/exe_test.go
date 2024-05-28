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
	if string(stdout) != "Hello, World!\n" {
		t.Fatalf("Unexpected stdout: %s", stdout)
	}
}

func TestExecuteAsyncHandlesCommandError(t *testing.T) {
	command := "invalid_command"
	execResult, err := ExecuteAsync(command)
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if execResult != nil {
		t.Fatalf("Expected nil ExecutionResult, got %v", execResult)
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
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if execResult != nil {
		t.Fatalf("Expected nil ExecutionResult, got %v", execResult)
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
	command := "/bin/bash -c '/usr/bin/sleep 5'"
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
	if stderr == "" {
		t.Fatalf("Expected stderr, got empty")
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
	command := "/bin/bash -c '/usr/bin/sleep 5'"
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
	command := "/bin/bash -c '/usr/bin/sleep 5'"
	execResult, err := ExecuteAsyncWithTimeout(command, 3*time.Second)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	<-execResult.Finished
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
	if err == nil {
		t.Fatalf("Expected error executing command, got nil")
	}
	if execResult != nil {
		t.Fatalf("Expected nil ExecutionResult, got %v", execResult)
	}
}

func TestExecuteAsyncWithTimeoutZeroTimeout(t *testing.T) {
	command := "/bin/bash -c '/usr/bin/sleep 5'"
	execResult, err := ExecuteAsyncWithTimeout(command, 0)
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	<-execResult.Finished
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
