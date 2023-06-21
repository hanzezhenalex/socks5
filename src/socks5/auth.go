package socks5

import (
	"context"
	"github.com/hanzezhenalex/socks5/src/auth"
	"net"
)

const (
	noAuth = uint8(0x00)

	authNoAuth = "noAuth"
)

type Authenticator interface {
	Name() string
	Method() uint8
	Handle(ctx context.Context, conn net.Conn) (auth.Info, error)
}

var emptyAuthInfo = auth.Info{}

type NoAuth struct{}

func (n NoAuth) Name() string {
	return authNoAuth
}

func (n NoAuth) Method() uint8 {
	return noAuth
}

func (n NoAuth) Handle(_ context.Context, _ net.Conn) (auth.Info, error) {
	return emptyAuthInfo, nil
}
