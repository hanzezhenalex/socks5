package socks5

import (
	"context"
	"net"

	"github.com/hanzezhenalex/socks5/src"
)

const (
	noAuth = uint8(0x00)

	authNoAuth = "noAuth"
)

type Authenticator interface {
	Name() string
	Method() uint8
	Handle(ctx context.Context, conn net.Conn) (src.AuthInfo, error)
}

var emptyAuthInfo = src.AuthInfo{}

type NoAuth struct{}

func (n NoAuth) Name() string {
	return authNoAuth
}

func (n NoAuth) Method() uint8 {
	return noAuth
}

func (n NoAuth) Handle(_ context.Context, _ net.Conn) (src.AuthInfo, error) {
	return emptyAuthInfo, nil
}
