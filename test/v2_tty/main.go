package main

import (
	"fmt"
	v2 "github.com/BGrewell/go-execute/v2"
)

func main() {
	e := v2.NewExecutorAsUser("whoopsie", []string{"BOB=YOUR_UNCLE"})
	err := e.ExecuteTTY("/bin/bash")
	if err != nil {
		fmt.Printf("failed to execute: %v\n", err)
	}
}
