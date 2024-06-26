package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"runtime"
)

func main() {

	if runtime.GOOS == "windows" {
		panic("windows not yet supported")
	}

	e := execute.NewExecutor(
		execute.WithUser("whoopsie"),
	)
	result, err := e.Execute("whoami")
	if err != nil {
		fmt.Printf("[-] Error: %v\n", err)
	}
	fmt.Println(result)
}
