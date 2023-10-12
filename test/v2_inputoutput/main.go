package main

import (
	v2 "github.com/BGrewell/go-execute/v2"
)

func main() {
	e := v2.NewExecutorAsUser("whoopsie", []string{"BOB=YOUR_UNCLE"})
	e.ExecuteTTY("/bin/bash")
}
