package socks5

import (
	"io"
	"net"
	"testing"

	"github.com/hanzezhenalex/socks5/src"
	"github.com/hanzezhenalex/socks5/src/socks5/protocol"

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

func createSocksServer() (*Server, error) {
	connMngr := &src.ConnectionManagement{}
	authMngr := &struct{}{}
	cfg := Config{
		IP:      "localhost",
		Port:    "8099",
		Command: []string{"connect"},
		Auth:    []string{"noAuth"},
	}
	return NewServer(cfg, connMngr, authMngr)
}

func Test_CommandConnect(t *testing.T) {
	rq := require.New(t)

	srv, err := createSocksServer()
	rq.NoError(err)
	defer func() {
		srv.Close()
	}()

	echoServer := TcpEchoServer{
		addr: "127.0.0.1:8098",
	}
	rq.NoError(echoServer.start())

	go func() {
		buf := make([]byte, 1024)

		conn, err := net.Dial("tcp", srv.config.Addr())
		rq.NoError(err)

		_, err = conn.Write([]byte{
			protocol.Socks5Version,
			0x01,
			protocol.Connect.Method,
		})
		rq.NoError(err)

		_, err = io.ReadFull(conn, buf[:2])
		rq.NoError(err)
		rq.Equal(uint8(protocol.Socks5Version), buf[0])
		rq.Equal(protocol.Connect.Method, buf[1])

		_, err = conn.Write([]byte{
			protocol.Socks5Version,
			protocol.Connect.Method,
			protocol.Rsv,
		})
		rq.NoError(err)
		addr := protocol.ParseAddr(echoServer.addr)
		rq.NotNil(addr)
		_, err = conn.Write(addr)
		rq.NoError(err)

		success, _, err := protocol.ReadCommandNegotiationReq(conn, buf)
		rq.NoError(err)
		rq.Equal(protocol.CommandNegoSucceed, success)

		testMsg := "hello socks5"
		_, err = conn.Write([]byte(testMsg))
		rq.NoError(err)

		data, err := io.ReadAll(conn)
		rq.NoError(err)
		rq.Equal(testMsg, string(data))
	}()

	rq.NoError(echoServer.onConnection())
}
