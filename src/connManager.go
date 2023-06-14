package src

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type ConnectionManager interface {
	Pipe(ctx context.Context, authInfo AuthInfo, from, to net.Conn) error
	DialTCP(ctx context.Context, authInfo AuthInfo, addr string) (net.Conn, net.Addr, error)
	Close()
}

type pipe struct {
	connMngr *ConnectionManagement
	uuid     uuid.UUID
	authInfo AuthInfo
	from, to net.Conn
}

func (p *pipe) copyTo(ch chan error, from, to net.Conn) {
	_, err := io.Copy(from, to)
	ch <- err
}

func (p *pipe) pipe(ctx context.Context, ch chan error) {
	received := 0
	connClosed := false
	closeConn := func() {
		if connClosed == true {
			return
		}
		_ = p.from.Close()
		_ = p.to.Close()
		connClosed = true
	}
	tracer := logrus.WithField("id", p.uuid)
	tracer.Info("start piping")

	go p.copyTo(ch, p.from, p.to)
	go p.copyTo(ch, p.to, p.from)

LOOP:
	for {
		select {
		case <-p.connMngr.stopCh:
			closeConn()
		case <-ctx.Done():
			closeConn()
		case err := <-ch:
			received++
			if err != nil && !(ReadOnClosedSocketError(err) && received == 2) {
				tracer.Errorf("an error happened when piping, err=%s", err.Error())
			}
			closeConn()
			if received == 2 {
				break LOOP
			}
		}
	}

	tracer.Info("pipe closed")
	p.connMngr.mutex.Lock()
	defer p.connMngr.mutex.Unlock()
	delete(p.connMngr.pipes, p.uuid)
	p.connMngr.wg.Done()
}

type ConnectionManagement struct {
	wg     sync.WaitGroup
	mutex  sync.Mutex
	pipes  map[uuid.UUID]*pipe
	closed bool
	stopCh chan struct{}
}

func NewConnectionManagement() *ConnectionManagement {
	return &ConnectionManagement{
		pipes:  make(map[uuid.UUID]*pipe),
		stopCh: make(chan struct{}),
	}
}

func (connMngr *ConnectionManagement) Pipe(ctx context.Context, authInfo AuthInfo, from, to net.Conn) error {
	newPipe := func() *pipe {
		p := &pipe{
			connMngr: connMngr,
			from:     from,
			to:       to,
			authInfo: authInfo,
			uuid:     GetIDFromContext(ctx),
		}

		connMngr.mutex.Lock()
		defer connMngr.mutex.Unlock()

		if connMngr.closed {
			return nil
		}
		connMngr.pipes[p.uuid] = p
		connMngr.wg.Add(1)
		return p
	}

	ch := make(chan error, 2)
	p := newPipe()

	if p == nil {
		return fmt.Errorf("connection manager has been closed")
	}
	go p.pipe(ctx, ch)
	return nil
}

func (connMngr *ConnectionManagement) DialTCP(_ context.Context, _ AuthInfo, addr string) (net.Conn, net.Addr, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.RemoteAddr(), nil
}

func (connMngr *ConnectionManagement) Close() {
	logrus.Info("start closing connection management")
	connMngr.mutex.Lock()
	if connMngr.closed == false {
		connMngr.closed = true
		close(connMngr.stopCh)
	} else {
		connMngr.mutex.Unlock()
		return
	}

	connMngr.mutex.Unlock()
	connMngr.wg.Wait()
	logrus.Info("connection management closed")
}

func ReadOnClosedSocketError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}
