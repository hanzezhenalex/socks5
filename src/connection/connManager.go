package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/hanzezhenalex/socks5/src"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	Pipe(ctx context.Context, authInfo src.AuthInfo, from, to net.Conn, target string) error
	DialTCP(ctx context.Context, authInfo src.AuthInfo, addr string) (net.Conn, net.Addr, error)
	Close()
	ListConnections(ctx context.Context, authInfo src.AuthInfo) ([]byte, error)
}

type pipe struct {
	connMngr *LocalManagement
	uuid     uuid.UUID
	authInfo src.AuthInfo
	from, to net.Conn
	target   string
}

func (p *pipe) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UUID   uuid.UUID `json:"uuid"`
		From   string    `json:"source"`
		Target string    `json:"target"`
	}{
		UUID:   p.uuid,
		From:   p.from.RemoteAddr().String(),
		Target: p.target,
	})
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

type LocalManagement struct {
	wg     sync.WaitGroup
	mutex  sync.Mutex
	pipes  map[uuid.UUID]*pipe
	closed bool
	stopCh chan struct{}
}

func NewConnectionManagement() *LocalManagement {
	return &LocalManagement{
		pipes:  make(map[uuid.UUID]*pipe),
		stopCh: make(chan struct{}),
	}
}

func (connMngr *LocalManagement) Pipe(ctx context.Context, authInfo src.AuthInfo, from, to net.Conn, target string) error {
	newPipe := func() *pipe {
		p := &pipe{
			connMngr: connMngr,
			from:     from,
			to:       to,
			authInfo: authInfo,
			uuid:     src.GetIDFromContext(ctx),
			target:   target,
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

func (connMngr *LocalManagement) DialTCP(_ context.Context, _ src.AuthInfo, addr string) (net.Conn, net.Addr, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.RemoteAddr(), nil
}

func (connMngr *LocalManagement) Close() {
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

func (connMngr *LocalManagement) ListConnections(ctx context.Context, authInfo src.AuthInfo) ([]byte, error) {
	connMngr.mutex.Lock()
	defer connMngr.mutex.Unlock()

	return json.Marshal(connMngr.pipes)
}

func ReadOnClosedSocketError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}
