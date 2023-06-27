package util

import "testing"

func TestWaitCloser(t *testing.T) {
	wc := NewWaitCloser()

	for i := 0; i < 10; i++ {
		go func() {
			_, ch := wc.Add()
			<-ch
			wc.Done()
		}()
	}

	wc.Close()
}
