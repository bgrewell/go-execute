package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"runtime"
)

func main() {

	e := execute.NewExecutor(
		execute.WithEnvironment([]string{"WHERE_AM_I=INSIDE_TTY"}),
	)

	cmd := "/bin/bash"
	if runtime.GOOS == "windows" {
		cmd = "cmd.exe"
	}

	err := e.ExecuteTTY(cmd)
	if err != nil {
		fmt.Printf("failed to execute: %v\n", err)
	}
}
