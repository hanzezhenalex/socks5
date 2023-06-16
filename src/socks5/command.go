package socks5

import (
	"context"
	"github.com/hanzezhenalex/socks5/src/connection"
	"net"
	"strings"

	"github.com/hanzezhenalex/socks5/src"
)

const (
	connect = uint8(0x01)

	cmdConnect = "connect"
)

type Commander interface {
	Name() string
	Method() uint8
	Handle(ctx context.Context, authInfo src.AuthInfo, target Addr, conn net.Conn, buf []byte) (net.Conn, error)
}

type ConnectCommandor struct {
	connMngr connection.Manager
}

func NewConnectCommandor(connMngr connection.Manager) ConnectCommandor {
	return ConnectCommandor{connMngr: connMngr}
}

func (c ConnectCommandor) Name() string {
	return cmdConnect
}

func (c ConnectCommandor) Method() uint8 {
	return connect
}

func (c ConnectCommandor) Handle(ctx context.Context, authInfo src.AuthInfo, target Addr, conn net.Conn, buf []byte) (net.Conn, error) {
	to, addr, err := c.connMngr.DialTCP(ctx, authInfo, target.String())
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "refused") {
			return nil, connectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			return nil, networkUnreachable
		}
		return nil, hostUnreachable
	}
	if err := writeCommandNegotiationReply(conn, buf, addr.String()); err != nil {
		_ = to.Close()
		return nil, err
	}
	return to, nil
}
