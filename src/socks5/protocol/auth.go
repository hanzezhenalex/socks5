package protocol

import (
	"context"
	"net"

	"github.com/hanzezhenalex/socks5/src"
)

var AuthHandlers = map[string]AuthHandler{
	NoAuth.Name: NoAuth,
}

type AuthHandler struct {
	Name   string
	Method byte
	Handle func(ctx context.Context, conn net.Conn, authMngr src.AuthManager) (src.AuthInfo, error)
}

var NoAuth = AuthHandler{
	Name:   "noAuth",
	Method: 0x00,
	Handle: func(_ context.Context, _ net.Conn, _ src.AuthManager) (src.AuthInfo, error) {
		return src.AuthInfo{}, nil
	},
}
