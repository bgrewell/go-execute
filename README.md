# go-execute
Go-Execute is a simple library that wraps some of the command execution functionality in go with the goal of making it
easier to use, especially in the context of running more complicated commands. It does this by abstracting some of the 
more complicated parts of the command execution process and providing a simple interface to run commands and get the
output with a variety of input and output options.

*Note: Version 1 of this library is still available to maintain compatibility with existing code. However, it is no
longer supported and will not receive any updates or bug fixes. It is recommended to upgrade to version 2. Examples
shown in this README are for version 2.*

## Installation

```bash
go get github.com/bgrewell/go-execute/v2
```

## Usage

### Simple Execution

```go
package main

import (
    "fmt"
    "github.com/bgrewell/go-execute/v2"
)

func main() {
    output, err := execute.Execute("whoami")
    if err != nil {
        fmt.Println("Error running command:", err)
        return
    }
    fmt.Println("Output:", output)
}
```

### Environment Variables Set Before Execution

```go
package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"runtime"
)

func main() {

	// Create a new executor with an env var set
	e := execute.NewExecutorWithEnv([]string{"BOB=YOUR_UNCLE"})

	// Use the appropriate command to print the enviornment variables
	command := "env"
	if runtime.GOOS == "windows" {
		command = "cmd /C set"
	}

	// Execute and print env
	result, err := e.Execute(command)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

```