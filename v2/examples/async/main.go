package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"github.com/bgrewell/go-execute/v2/pkg"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {

	logger := pkg.NewZapLogger()
	logger.SetLevel(zapcore.DebugLevel)
	execute.SetLogger(logger)

	ex := execute.NewExecutor()

	cmd := "/bin/bash -c for i in {10..1}; do echo -ne '$i\\033[0K\\r'; sleep 1; done; echo 'DONE'\n"
	if runtime.GOOS == "windows" {
		cmd = "cmd.exe /c for /l %i in (10,-1,1) do (echo %i & timeout /t 1 /nobreak >nul) & echo DONE"
	}
	cmd = "whoami"

	result, err := ex.ExecuteAsync(cmd)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		wg.Done()
		readAndOutput(result.Stdout, os.Stdout, result.Finished)
	}()

	go func() {
		wg.Done()
		readAndOutput(result.Stderr, os.Stderr, result.Finished)
	}()

	<-result.Finished
	wg.Wait()

}

func readAndOutput(r io.Reader, w io.Writer, finished <-chan error) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-finished:
			fmt.Println("REMOVE ME: Finished reading")
			return
		default:
			n, err := r.Read(buf)
			if err != nil {
				if err != io.EOF {
					// Handle read error if necessary
					fmt.Printf("Error reading from stream: %v\n", err)
				}
				return
			}
			if n > 0 {
				fmt.Println("REMOVE ME: Writing to stream")
				w.Write(buf[:n])
			}
		}
	}
}
