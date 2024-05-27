//go:build integration
// +build integration

package execute

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	result, err := Execute("whoami")
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}
	t.Logf("Result: %v", result)
}

func TestExecuteSeparate(t *testing.T) {
	command := "whoami"
	stdout, stderr, err := ExecuteSeparate(command)
	if err != nil {
		t.Fatalf("Error reading stderr: %v", err)
	} else if len(stderr) > 0 {
		t.Fatalf("Unexpected stderr: %s", stderr)

	}
	t.Logf("%s: %s", command, stdout)
}

func TestExecuteAsync(t *testing.T) {
	command := "whoami"
	execResult, err := ExecuteAsync(command)
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
	} else if len(stderr) > 0 {
		t.Fatalf("Unexpected stderr: %s", stderr)

	}
	t.Logf("%s: %s", command, stdout)
}

func TestExecuteAsyncWithInput(t *testing.T) {
	command := "cat"

	input := io.NopCloser(strings.NewReader("testing"))
	execResult, err := ExecuteAsyncWithInput(command, input)
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
	} else if len(stderr) > 0 {
		t.Fatalf("Unexpected stderr: %s", stderr)

	}
	t.Logf("%s: %s", command, stdout)
}

func TestExecuteAsyncWithTimeoutSuccess(t *testing.T) {
	command := "/bin/bash -c '/usr/bin/sleep 5'"
	execResult, err := ExecuteAsyncWithTimeout(command, 6*time.Second)
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
	} else if len(stderr) > 0 {
		t.Fatalf("Unexpected stderr: %s", stderr)

	}
	t.Logf("%s: %s", command, stdout)
}
