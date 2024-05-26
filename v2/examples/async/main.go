package main

import (
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {

	ex := execute.NewExecutor()
	wg := sync.WaitGroup{}

	cmd := "/bin/bash -c for i in {10..1}; do echo -ne '$i\\033[0K\\r'; sleep 1; done; echo 'DONE'\n"
	if runtime.GOOS == "windows" {
		cmd = "cmd.exe /c for /l %i in (10,-1,1) do (echo %i & timeout /t 1 /nobreak >nul) & echo DONE"
	}
	fmt.Println("Running command: ", cmd)

	result, err := ex.ExecuteAsync(cmd)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}

	wg.Add(1)
	go func() {
		readAndOutput(result.Stdout, os.Stdout, result.Finished, &wg)
	}()

	wg.Add(1)
	go func() {
		readAndOutput(result.Stderr, os.Stderr, result.Finished, &wg)
	}()

	<-result.Finished
	fmt.Println("REMOVE ME: Finished")
	wg.Wait()
	fmt.Println("REMOVE ME: Done")
}

func readAndOutput(r io.Reader, w io.Writer, finished <-chan error, wg *sync.WaitGroup) {
	buf := make([]byte, 1024)
	for {
		select {
		case <-finished:
			fmt.Println("REMOVE ME: Finished reading")
			wg.Done()
			return
		default:
			n, err := r.Read(buf)
			if err != nil {
				if err != io.EOF {
					// Handle read error if necessary
					fmt.Printf("Error reading from stream: %v\n", err)
					wg.Done()
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
