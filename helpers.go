package execute

import (
	"bufio"
	"errors"
	"io"
	"time"
)

// SignalRead checks if the reader is ready to be read from before the timeout.
func SignalRead(reader io.ReadCloser, timeout time.Duration) (ready bool, err error) {
	// Channels to signal when data is ready or an error occurred
	readyChan := make(chan bool)
	errChan := make(chan error)

	// Goroutine to perform the reading
	go func() {
		bufferedReader := bufio.NewReader(reader)
		_, err := bufferedReader.Peek(1) // Just peeking the first byte
		if err != nil {
			if err != bufio.ErrBufferFull {
				errChan <- err
				return
			}
		}
		readyChan <- true
	}()

	// Select to wait on channels or timeout
	select {
	case <-readyChan:
		// Data is ready
		return true, nil
	case err := <-errChan:
		// An error occurred while reading
		return false, err
	case <-time.After(timeout):
		// Timeout occurred
		return false, errors.New("read timeout")
	}
}
