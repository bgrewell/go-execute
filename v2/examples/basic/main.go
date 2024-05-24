package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
)

func main() {

	// Create a new executor
	ex := execute.NewExecutor()

	// Run a basic command
	result, err := ex.Execute("whoami")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
