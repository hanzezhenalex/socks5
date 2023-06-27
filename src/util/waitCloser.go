package util

import "sync"

type status string

const (
	running status = "running"
	closing status = "closing"
	closed  status = "closed"
)

type WaitCloser struct {
	lock     sync.Mutex
	status   status
	wg       sync.WaitGroup
	stopCh   chan struct{}
	closedCh chan struct{}
}

func NewWaitCloser() *WaitCloser {
	return &WaitCloser{
		status:   running,
		stopCh:   make(chan struct{}),
		closedCh: make(chan struct{}),
	}
}

func (wc *WaitCloser) Add() (bool, chan struct{}) {
	wc.lock.Lock()
	defer wc.lock.Unlock()

	if wc.status == running {
		wc.wg.Add(1)
		return true, wc.stopCh
	}
	return false, nil
}

func (wc *WaitCloser) Done() {
	wc.wg.Done()
}

func (wc *WaitCloser) Close() {
	wc.lock.Lock()

	switch wc.status {
	case running:
		wc.status = closing
		close(wc.stopCh)

		wc.lock.Unlock()

		wc.wg.Wait()
		close(wc.closedCh)

		wc.lock.Lock()
		defer wc.lock.Unlock()

		wc.status = closed
	case closing:
		wc.lock.Unlock()
		<-wc.closedCh
	default:
		wc.lock.Unlock()
	}
}
