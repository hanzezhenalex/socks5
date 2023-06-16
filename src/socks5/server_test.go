package socks5

import (
	"github.com/hanzezhenalex/socks5/src/connection"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

type TcpEchoServer struct {
	addr string
	l    net.Listener
}

func (srv *TcpEchoServer) start() error {
	l, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return err
	}
	srv.l = l
	return nil
}

func (srv *TcpEchoServer) onConnection() error {
	defer func() {
		_ = srv.l.Close()
	}()
	conn, err := srv.l.Accept()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)

	nr, err := conn.Read(buf)
	if err != nil {
		return err
	}
	nw, err := conn.Write(buf[:nr])
	if err != nil {
		return err
	}
	if nw != nr {
		return err
	}

	return conn.Close()
}

func createSocksServer() (*Server, error, chan error) {
	connMngr := connection.NewConnectionManagement()
	authMngr := &struct{}{}
	cfg := Config{
		IP:      "localhost",
		Port:    "8099",
		Command: []string{"connect"},
		Auth:    []string{"noAuth"},
	}
	ch := make(chan error, 1)
	srv, err := NewServer(cfg, connMngr, authMngr)
	go func() {
		ch <- srv.Start()
	}()
	return srv, err, ch
}

func Test_CommandConnect(t *testing.T) {
	rq := require.New(t)

	srv, err, ch := createSocksServer()
	rq.NoError(err)

	echoServer := TcpEchoServer{
		addr: "127.0.0.1:8098",
	}
	rq.NoError(echoServer.start())

	go func() {
		buf := make([]byte, 1024)

		conn, err := net.Dial("tcp", srv.config.Addr())
		rq.NoError(err)

		_, err = conn.Write([]byte{
			version,
			0x01,
			noAuth,
		})
		rq.NoError(err)

		_, err = io.ReadFull(conn, buf[:2])
		rq.NoError(err)
		rq.Equal(version, buf[0])
		rq.Equal(noAuth, buf[1])

		_, err = conn.Write([]byte{
			version,
			connect,
			rsv,
		})
		rq.NoError(err)
		addr := ParseAddr(echoServer.addr)
		rq.NotNil(addr)
		_, err = conn.Write(addr)
		rq.NoError(err)

		success, _, err := readCommandNegotiationReq(conn, buf)
		rq.NoError(err)
		rq.Equal(commandNegoSucceed, success)

		testMsg := "hello socks5"
		_, err = conn.Write([]byte(testMsg))
		rq.NoError(err)

		data, err := io.ReadAll(conn)
		rq.NoError(err)
		rq.Equal(testMsg, string(data))
	}()

	rq.NoError(echoServer.onConnection())
	srv.Close()
	<-ch
	close(ch)
}
