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

	shell := "/bin/bash"
	if runtime.GOOS == "windows" {
		shell = "powershell"
	}

	ex := execute.NewExecutor(
		execute.WithShell(shell),
	)

	cmd := "ls -laR /usr/share"
	if runtime.GOOS == "windows" {
		cmd = "cmd.exe /c for /l %i in (10,-1,1) do (echo %i & timeout /t 1 /nobreak >nul) & echo DONE"
	}

	result, err := ex.ExecuteAsync(cmd)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		readAndOutput(result.Stdout, os.Stdout, result.Finished)
		wg.Done()
	}()

	go func() {
		readAndOutput(result.Stderr, os.Stderr, result.Finished)
		wg.Done()
	}()

	<-result.Finished
	wg.Wait()

}

func readAndOutput(r io.Reader, w io.Writer, finished <-chan error) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-finished:
			n, err := r.Read(buf)
			if err != nil {
				if err != io.EOF {
					// Handle read error if necessary
					fmt.Printf("Error reading from stream: %v\n", err)
				}
			}
			if n > 0 {
				_, err = w.Write(buf[:n])
				if err != nil {
					fmt.Printf("Error writing to stream: %v\n", err)
				}
			}
			return
		default:
			n, err := r.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Error reading from stream: %v\n", err)
				}
				return
			}
			if n > 0 {
				_, err = w.Write(buf[:n])
				if err != nil {
					fmt.Printf("Error writing to stream: %v\n", err)
				}
			}
		}
	}
}
