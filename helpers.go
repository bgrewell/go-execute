package execute

import (
	"bufio"
	"errors"
	"time"
)

// SignalRead checks if the reader is ready to be read from before the timeout.
func SignalRead(bufferedReader *bufio.Reader, timeout time.Duration) (ready bool, err error) {
	readyChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		_, err := bufferedReader.Peek(1)
		if err != nil {
			if err != bufio.ErrBufferFull {
				errChan <- err
				return
			}
		}
		readyChan <- true
	}()

	select {
	case <-readyChan:
		return true, nil
	case err := <-errChan:
		return false, err
	case <-time.After(timeout):
		return false, errors.New("read timeout")
	}
}
