package connection

import (
	"context"
	"fmt"
	"github.com/hanzezhenalex/socks5/src/auth"
	"github.com/hanzezhenalex/socks5/src/util"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/hanzezhenalex/socks5/src"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Manager interface {
	Pipe(ctx context.Context, authInfo auth.Info, from, to net.Conn, target string) error
	DialTCP(ctx context.Context, authInfo auth.Info, addr string) (net.Conn, net.Addr, error)
	Close()
	ListConnections(ctx context.Context, authInfo auth.Info) []PipeInfo
}

type PipeInfo struct {
	UUID   uuid.UUID `json:"uuid"`
	From   string    `json:"source"`
	Target string    `json:"target"`
}

type pipe struct {
	connMngr *LocalManagement
	uuid     uuid.UUID
	authInfo auth.Info
	from, to net.Conn
	target   string
	stopCh   chan struct{}
}

func (p *pipe) info() PipeInfo {
	return PipeInfo{
		UUID:   p.uuid,
		From:   p.from.RemoteAddr().String(),
		Target: p.target,
	}
}

func (p *pipe) copyTo(ch chan error, from, to net.Conn) {
	_, err := io.Copy(from, to)
	ch <- err
}

func (p *pipe) pipe(ctx context.Context, ch chan error) {
	tracer := logrus.WithField("id", p.uuid)
	tracer.Info("start piping")

	defer func() {
		tracer.Info("pipe closed")
		p.connMngr.pipes.Delete(p.uuid)
		p.connMngr.closer.Done()
	}()

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

	go p.copyTo(ch, p.from, p.to)
	go p.copyTo(ch, p.to, p.from)

LOOP:
	for {
		select {
		case <-p.stopCh:
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
}

type LocalManagement struct {
	pipes  sync.Map
	closer *util.WaitCloser
}

func NewConnectionManagement() *LocalManagement {
	return &LocalManagement{
		closer: util.NewWaitCloser(),
	}
}

func (connMngr *LocalManagement) Pipe(ctx context.Context, authInfo auth.Info, from, to net.Conn, target string) error {
	newPipe := func() *pipe {
		ok, ch := connMngr.closer.Add()
		if !ok {
			return nil
		}

		p := &pipe{
			connMngr: connMngr,
			from:     from,
			to:       to,
			authInfo: authInfo,
			uuid:     src.GetIDFromContext(ctx),
			target:   target,
			stopCh:   ch,
		}

		connMngr.pipes.Store(p.uuid, p)
		return p
	}

	ch := make(chan error, 2)
	p := newPipe()

	if p == nil {
		return fmt.Errorf("connection manager has been closed")
	}
	logrus.Debugf("pipe created, target=%s", p.target)

	go p.pipe(ctx, ch)
	return nil
}

func (connMngr *LocalManagement) DialTCP(_ context.Context, _ auth.Info, addr string) (net.Conn, net.Addr, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	return conn, conn.RemoteAddr(), nil
}

func (connMngr *LocalManagement) Close() {
	logrus.Info("start closing connection management")

	connMngr.closer.Close()

	logrus.Info("connection management closed")
}

func (connMngr *LocalManagement) ListConnections(ctx context.Context, authInfo auth.Info) []PipeInfo {
	var infos []PipeInfo

	connMngr.pipes.Range(func(_, p any) bool {
		infos = append(infos, p.(*pipe).info())
		return true
	})

	return infos
}

func ReadOnClosedSocketError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}
