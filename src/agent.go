package src

import (
	"context"
	"fmt"
	"github.com/hanzezhenalex/socks5/src/protocol"
	"net"

	"github.com/hanzezhenalex/socks5/src/protocol/auth"
	"github.com/hanzezhenalex/socks5/src/protocol/cmd"
)

type Mode string

const (
	local  Mode = "local"
	remote Mode = "remote"
)

type ClientConfig struct {
	LocalIp    string
	LocalPort  int
	ServerIp   string
	ServerPort int
	Mode       Mode
}

type Client struct {
	cfg             ClientConfig
	commandHandlers map[byte]cmd.Handler
	authHandlers    map[byte]auth.Handler

	listener net.Listener
}

func (c *Client) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.cfg.LocalIp, c.cfg.LocalPort))

	if err != nil {
		return err
	}

	c.listener = listener

	for {
		conn, err := c.listener.Accept()
		if err != nil {

		}
		go c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	methods, err := protocol.ReadMethodNegotiationReq(conn)
	if err != nil {
		return
	}

	authHandler := c.selectAuthMethodHandler(methods)
	if authHandler == nil {
		if _, err = conn.Write(protocol.NoAcceptableMethods); err != nil {

		}
		return
	}

	if err = authHandler.Handle(context.Background(), conn); err != nil {

		return
	}

}

func (c *Client) selectAuthMethodHandler(supported []byte) auth.Handler {
	for _, m := range supported {
		if handler, ok := c.authHandlers[m]; ok {
			return handler
		}
	}
	return nil
}
