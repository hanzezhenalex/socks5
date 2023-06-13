package socks5

import (
	"context"
	"fmt"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/hanzezhenalex/socks5/src"
	"github.com/hanzezhenalex/socks5/src/socks5/protocol"
)

type Pipe struct {
	id       uuid.UUID
	srv      *Server
	from     net.Conn
	to       net.Conn
	tracer   *logrus.Entry
	authInfo src.AuthInfo
	command  string
	target   protocol.Addr
	cancel   context.CancelFunc
	isPiping bool
}

func (pipe *Pipe) run() {
	ctx, cancel := context.WithCancel(context.Background())
	pipe.cancel = cancel

	if err := pipe.handshake(ctx); err != nil {
		if _, ok := err.(protocol.NetworkError); !ok {
			_ = pipe.from.Close()
		}
		if socksErr, ok := err.(protocol.SocksError); ok {
			socksErr.SendErrorReply(pipe.from)
		}
		pipe.tracer.Errorf("an error happens when handshake, err=%s", err.Error())
		return
	}

	pipe.tracer.Infof("handshake successfully, piping now, target=%s", pipe.target.String())

	pipe.isPiping = true
	if err := pipe.srv.connMngr.Proxy(ctx, pipe.authInfo, pipe.from, pipe.to); err != nil {
		pipe.tracer.Errorf("an error happens when piping, err=%s", err.Error())
	}
}

func (pipe *Pipe) handshake(ctx context.Context) error {
	var (
		buf = make([]byte, protocol.MaxAddrLen)
		err error
	)

	methods, err := protocol.ReadMethodNegotiationReq(pipe.from, buf)
	if err != nil {
		return fmt.Errorf("fail to read method negotiation req, err=%w", err)
	}

	authHandler, err := pipe.srv.selectAuthMethod(methods)
	if err != nil {
		return err
	}

	if err := protocol.WriteMethodNegotiationReply(authHandler.Method, pipe.from); err != nil {
		return fmt.Errorf("fail to write method negotiation reply, err=%w", err)
	}

	info, err := authHandler.Handle(ctx, pipe.from, pipe.srv.authMngr)
	if err != nil {
		return err
	}
	pipe.authInfo = info

	cmd, targetAddr, err := protocol.ReadCommandNegotiationReq(pipe.from, buf)
	if err != nil {
		return fmt.Errorf("fail to read command negotiation req, err=%w", err)
	}
	pipe.target = targetAddr

	commandHandler, err := pipe.srv.getCommandMethod(cmd)
	if err != nil {
		return err
	}
	pipe.command = commandHandler.Name

	to, addr, err := commandHandler.Handle(ctx, pipe.authInfo, pipe.target, pipe.from, pipe.srv.connMngr)
	if err != nil {
		return err
	}
	pipe.to = to
	if err := protocol.WriteCommandNegotiationReply(pipe.from, buf, addr.String()); err != nil {
		_ = pipe.to.Close()
		return err
	}
	return nil
}
