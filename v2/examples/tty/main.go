package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"runtime"
)

func main() {

	var e execute.Executor

	if runtime.GOOS == "windows" {
		e = execute.NewExecutorWithEnv([]string{"WHERE_AM_I=INSIDE_TTY"})
		err := e.ExecuteTTY("cmd.exe")
		if err != nil {
			fmt.Printf("failed to execute: %v\n", err)
		}
	} else {
		e = execute.NewExecutorAsUser("whoopsie", []string{"WHERE_AM_I=INSIDE_TTY"})
		err := e.ExecuteTTY("/bin/bash")
		if err != nil {
			fmt.Printf("failed to execute: %v\n", err)
		}
	}

}
