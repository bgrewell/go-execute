package utilities

import (
	"fmt"
	"io"
	"sync"
)

// CopyAndClose copies the contents of the provided io.ReadCloser to the provided bytes.Buffer and closes the io.ReadCloser.
func CopyAndClose(buf io.Writer, r io.ReadCloser, ready *sync.WaitGroup, done chan struct{}) {
	defer r.Close()
	ready.Done()
	fmt.Println("signaled ready; beginning copy")
	_, err := io.Copy(buf, r)
	if err != nil {
		fmt.Printf("Error copying data: %v\n", err)
		return
	}
	fmt.Println("copy complete")
	done <- struct{}{}
}
