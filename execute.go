package execute

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/BGrewell/go-execute/internal/utilities"
	"io"
	"os"
	"os/exec"
	"time"
)

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	var output, stderr bytes.Buffer
	for i, cmd := range cmds {
		if i < len(cmds)-1 {
			if cmds[i+1].Stdin, pipeLineError = cmd.StdoutPipe(); pipeLineError != nil {
				return nil, nil, pipeLineError
			}
		} else {
			cmd.Stdout = &output
		}
		cmd.Stderr = &stderr
		if pipeLineError = cmd.Start(); pipeLineError != nil {
			return output.Bytes(), stderr.Bytes(), pipeLineError
		}
	}

	for _, cmd := range cmds {
		if pipeLineError = cmd.Wait(); pipeLineError != nil {
			return output.Bytes(), stderr.Bytes(), pipeLineError
		}
	}

	return output.Bytes(), stderr.Bytes(), nil
}

func ExecutePipedCmds(commands []string) (output string, err error) {
	if len(commands) < 2 {
		return "", fmt.Errorf("you must pass 2 or more commands to pipe them")
	}

	cmds := make([]*exec.Cmd, len(commands))
	for idx, command := range commands {
		cmds[idx], err = prepareCommand(command)
		if err != nil {
			return "", err
		}
	}

	bytesout, stderr, err := Pipeline(cmds...)
	if err != nil {
		return string(bytesout) + "\n" + string(stderr), err
	}
	if len(stderr) > 0 && string(stderr) != "" {
		return "", fmt.Errorf(string(stderr))
	}
	return string(bytesout), err
}

// ExecuteCmds executes a slice of commands and returns the output and errors from each
func ExecuteCmds(commands []string) (outputs []string, errs []error) {
	outputs = make([]string, len(commands))
	errs = make([]error, len(commands))
	for idx, command := range commands {
		outputs[idx], errs[idx] = ExecuteCmd(command)
	}
	return outputs, errs
}

// ExecuteCmd executes commands and returns the output and any errors
func ExecuteCmd(command string) (output string, err error) {
	exe, err := prepareCommand(command)
	if err != nil {
		return "", err
	}
	out, err := exe.CombinedOutput()
	return string(out), err
}

// ExecuteCmdEx executes commands and returns the stdout and stderr as separate strings
func ExecuteCmdEx(command string) (stdout string, stderr string, err error) {
	var bout, berr bytes.Buffer
	exe, err := prepareCommand(command)
	if err != nil {
		return "", "", err
	}
	exe.Stdout = &bout
	exe.Stderr = &berr
	err = exe.Run()
	return string(bout.Bytes()), string(berr.Bytes()), err
}

// ExecuteCmdWithEnvVars executes a command with the passed in env vars set and returns the results
func ExecuteCmdWithEnvVars(command string, vars []string) (stdout string, stderr string, err error) {
	var bout, berr bytes.Buffer
	exe, err := prepareCommand(command)
	if err != nil {
		return "", "", err
	}
	exe.Env = os.Environ()
	exe.Env = append(exe.Env, vars...)
	exe.Stdout = &bout
	exe.Stderr = &berr
	err = exe.Run()
	return string(bout.Bytes()), string(berr.Bytes()), err
}

// ExecuteCmdWithTimeout executes commands with a timeout. If the timeout occurs the command is terminated and an error is returned
func ExecuteCmdWithTimeout(command string, seconds int) (output string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()

	exe, err := prepareCommand(command)
	if err != nil {
		return "", err
	}

	outBytes, err := exec.CommandContext(ctx, exe.Path, exe.Args[1:]...).Output()
	if ctx.Err() == context.DeadlineExceeded {
		err = errors.New("command execution timeout exceeded")
	}
	return string(outBytes), err
}

// ExecutePowershell executes a command using powershell and returns the stdout, stderr and any error code
func ExecutePowershell(command string) (stdout string, stderr string, err error) {
	encCommand := conversions.ConvertToUTF16LEBase64String(command)

	var bout, berr bytes.Buffer
	exe := exec.Command("powershell.exe", "-NoProfile", "-enc", encCommand)
	exe.Stdout = &bout
	exe.Stderr = &berr
	err = exe.Run()

	return string(bout.Bytes()), string(berr.Bytes()), err
}

func ExecuteAsync(command string, env *[]string) (outPipe io.ReadCloser, errPipe io.ReadCloser, exitCode chan int, err error) {
	stdout, stderr, exitCode, _, err := ExecuteAsyncWithCancel(command, env)
	return stdout, stderr, exitCode, err
}

func ExecuteAsyncWithCancel(command string, env *[]string) (stdOut io.ReadCloser, stdErr io.ReadCloser, exitCode chan int, cancelToken context.CancelFunc, err error) {
	exitCode = make(chan int)
	ctx, cancel := context.WithCancel(context.Background())

	cmd, err := prepareCommand(command)
	if err != nil {
		cancel()
		return nil, nil, nil, nil, err
	}
	exe := exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	exe.Env = os.Environ()
	if env != nil {
		exe.Env = append(exe.Env, *env...)
	}
	stdOut, err = exe.StdoutPipe()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	stdErr, err = exe.StderrPipe()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	err = exe.Start()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutDone := make(chan struct{})
	stderrDone := make(chan struct{})

	go copyAndClose(stdoutDone, &stdoutBuf, stdOut)
	go copyAndClose(stderrDone, &stderrBuf, stdErr)

	// Wait for the command to complete or for the timeout
	done := make(chan error, 1)
	go func() {
		done <- exe.Wait()
	}()

	// Wait for completion or timeout
	select {
	case <-ctx.Done():
		defer cancel()
		return nil, nil, exitCode, nil, ctx.Err()
	case err := <-done:
		return io.NopCloser(bytes.NewReader(stdoutBuf.Bytes())), io.NopCloser(bytes.NewReader(stderrBuf.Bytes())), exitCode, cancel, err
	}
}

func prepareCommand(command string) (*exec.Cmd, error) {

	cmdParts, err := utilities.Fields(command) // Assuming utilities.Fields breaks the command string into parts
	if err != nil {
		return nil, err
	}

	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return nil, err
	}

	exe := exec.Command(binary, cmdParts[1:]...)
	return exe, nil
}

func copyAndClose(done chan struct{}, buf *bytes.Buffer, r io.ReadCloser) {
	defer close(done)
	defer r.Close()
	io.Copy(buf, r)
}
