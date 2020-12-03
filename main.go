package execute

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
	"github.com/BGrewell/go-conversions"
)

var (
	Debug = false
)

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

func ExecutePipedCmds(commands []string) (output string, err error) {
	if len(commands) < 2 {
		return "", fmt.Errorf("you must pass 2 or more commands to pipe them")
	}
	cmds := make([]*exec.Cmd, len(commands))
	for idx := 0; idx < len(commands); idx++ {
		// break the command into it's fields
		fmt.Println("cmd: " + commands[idx])
		cmdParts := strings.Fields(commands[idx])
		exe, err := exec.LookPath(cmdParts[0])
		if err != nil {
			return "", err
		}
		// setup execution
		cmd := exec.Command(exe, cmdParts[1:]...)
		cmds[idx] = cmd
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
	cmdParts := strings.Split(command, " ")
	exename, err := exec.LookPath(cmdParts[0])
	exe := exec.Command(exename, cmdParts[1:]...)
	out, err := exe.CombinedOutput()
	return string(out), err
}

// ExecuteCmdEx executes commands and returns the stdout and stderr as seperate strings
func ExecuteCmdEx(command string) (stdout string, stderr string, err error) {
	var bout, berr bytes.Buffer
	cmdParts := strings.Split(command, " ")
	exename, err := exec.LookPath(cmdParts[0])
	exe := exec.Command(exename, cmdParts[1:]...)
	exe.Stdout = &bout
	exe.Stderr = &berr
	err = exe.Run()
	return string(bout.Bytes()), string(berr.Bytes()), err
}

// ExecuteCmdWithTimeout executes commands with a timeout. If the timeout occurs the command is terminated and an error is returned
func ExecuteCmdWithTimeout(command string, seconds int) (output string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel()
	cmdParts := strings.Split(command, " ")
	exename, err := exec.LookPath(cmdParts[0])
	outBytes, err := exec.CommandContext(ctx, exename, cmdParts[1:]...).Output()
	if ctx.Err() == context.DeadlineExceeded {
		err = errors.New("command execution timeout exceeded")
	}
	return string(outBytes), err
}

// ExecutePowershell executes a command using powershell and returns the stdout, stderr and any error code
func ExecutePowershell(command string) (stdout string, stderr string, err error) {
	//command = strings.ReplaceAll(command,"\"", "\\\"")
	//command = strings.ReplaceAll(command, "'", "\\'")
	//command = fmt.Sprintf("'%s'", command)
	encCommand := conversions.ConvertToUTF16LEString(command)
	if Debug {
		fmt.Println(encCommand)
	}
	var bout, berr bytes.Buffer
	exename, err := exec.LookPath("powershell.exe")
	exe := exec.Command(exename, "-enc", encCommand)
	exe.Stdout = &bout
	exe.Stderr = &berr
	err = exe.Run()
	return string(bout.Bytes()), string(berr.Bytes()), err
}