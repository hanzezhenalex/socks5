package protocol

import (
	"context"
	"net"

	"github.com/hanzezhenalex/socks5/src"
)

var Connect = CommandHandler{
	Name:   "connect",
	Method: 0x01,
	Handle: func(ctx context.Context, authInfo src.AuthInfo, addr Addr, conn net.Conn,
		connMngr src.ConnectionManager) (net.Addr, error) {
		return connMngr.Proxy(ctx, authInfo, conn, addr.String())
	},
}

type CommandHandler interface {
	Name() string
	Method() uint8
	Handle(ctx context.Context, authInfo src.AuthInfo, addr Addr, conn net.Conn,
		connMngr src.ConnectionManager) (net.Addr, error)
}
