package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
)

func main() {

	// Just run a command without creating an executor
	result, err := execute.Execute("whoami")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
