package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
)

func main() {
	// Create a new PowerShell executor
	psExecutor := execute.NewPowerShellExecutor()

	// Execute a single command
	result, err := psExecutor.Execute("Get-Process")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Execute a script with parameters
	script := `param([string]$name) Write-Output "Hello, $name"`
	params := map[string]string{"name": "World"}
	scriptResult, err := psExecutor.ExecuteScript(script, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(scriptResult)

	// Start a PowerShell session
	err = psExecutor.StartSession()
	if err != nil {
		panic(err)
	}

	// Execute commands in the session
	sessionResult, err := psExecutor.ExecuteInSession("Get-Date")
	if err != nil {
		panic(err)
	}
	fmt.Println(sessionResult)
}
