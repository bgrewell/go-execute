package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"runtime"
)

func main() {

	// Create a new executor with an env var set
	e := execute.NewExecutor(
		execute.WithDefaultShell(),
		execute.WithEnvironment([]string{"BOB=YOUR_UNCLE"}),
	)

	// Use the appropriate command to print the enviornment variables
	command := "env"
	if runtime.GOOS == "windows" {
		command = "set"
	}

	// Execute and print env
	result, err := e.Execute(command)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
