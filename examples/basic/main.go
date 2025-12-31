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

	// Just run a command without creating an executor
	result, err := execute.Execute("whoami")
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
