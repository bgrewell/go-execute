package main

import (
	"fmt"
	v2 "github.com/BGrewell/go-execute/v2"
)

func main() {
	e := v2.NewExecutorAsUser("whoopsie", nil)
	result, err := e.Execute("whoami")
	if err != nil {
		fmt.Printf("[-] Error: %v\n", err)
	}
	fmt.Println(result)
}
