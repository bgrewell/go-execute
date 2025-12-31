package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"github.com/bgrewell/go-execute/v2/pkg"
	"go.uber.org/zap/zapcore"
)

func main() {

	// Not needed for this example, but useful for debugging
	logger := pkg.NewZapLogger()
	logger.SetLevel(zapcore.DebugLevel)
	execute.SetLogger(logger)

	// Create a new PowerShell executor
	e := execute.NewExecutor(
		execute.WithDefaultShell(),
	)

	// Execute a script with parameters
	script := `param([string]$name) Write-Output "Hello, $name"`
	params := map[string]string{"name": "World"}
	stdout, stderr, err := e.ExecuteScriptFromString(execute.ScriptTypePowerShell, script, nil, params)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
	fmt.Printf("Results: %s\n", stdout)
	if stderr != "" {
		fmt.Printf("Error: %s\n", stderr)
	}

}
