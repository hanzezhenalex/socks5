package cmd

import (
	"context"
	"net"

	socks5 "github.com/hanzezhenalex/socks5/src"
)

type Handler interface {
	Name() string
	Method() byte
	Handle(ctx context.Context, conn net.Conn, manager socks5.Proxier) error
}
