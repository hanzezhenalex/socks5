package src

import (
	"context"
	"net"
)

type ConnectionManager interface {
	Proxy(ctx context.Context, authInfo AuthInfo, from net.Conn, addr string) (net.Addr, error)
}
