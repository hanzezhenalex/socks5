package socks5

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/hanzezhenalex/socks5/src"
	"github.com/hanzezhenalex/socks5/src/socks5/protocol"
)

type Config struct {
	IP      string
	Port    string
	Auth    []string
	Command []string
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.IP, c.Port)
}

type Server struct {
	connMngr        src.ConnectionManager
	authMngr        src.AuthManager
	listener        net.Listener
	config          Config
	handlerLock     sync.Mutex
	commandHandlers map[byte]protocol.CommandHandler
	authHandlers    map[byte]protocol.AuthHandler
}

func NewServer(cfg Config, connMngr src.ConnectionManager, authMngr src.AuthManager) (*Server, error) {
	srv := &Server{
		config:          cfg,
		connMngr:        connMngr,
		authMngr:        authMngr,
		commandHandlers: make(map[byte]protocol.CommandHandler),
		authHandlers:    make(map[byte]protocol.AuthHandler),
	}

	logrus.Infof("register auth handlers, handlers=[%s]", strings.Join(cfg.Auth, ","))
	for _, name := range cfg.Auth {
		if err := srv.AddAuthHandler(name); err != nil {
			return nil, err
		}
	}

	logrus.Infof("register command handlers, handlers=[%s]", strings.Join(cfg.Command, ","))
	for _, name := range cfg.Command {
		if err := srv.AddCommandHandler(name); err != nil {
			return nil, err
		}
	}

	if err := srv.startTCPServer(); err != nil {
		return nil, fmt.Errorf("fail to start tcp server: %w", err)
	}

	return srv, nil
}

func (srv *Server) startTCPServer() error {
	logrus.Infof("start socks server, commands=%s", strings.Join(srv.config.Command, ","))
	l, err := net.Listen("tcp", srv.config.Addr())
	if err != nil {
		return err
	}
	srv.listener = l

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				logrus.Warningf("fail to accept tcp conn, err=%s", err.Error())
				return
			}
			go srv.handleConn(conn)
		}
	}()

	return nil
}

func (srv *Server) handleConn(conn net.Conn) {
	ctx := context.Background()

	var (
		authInfo src.AuthInfo
	)

	handshake := func() error {
		buf := make([]byte, protocol.MaxAddrLen)

		methods, err := protocol.ReadMethodNegotiationReq(conn, buf)
		if err != nil {
			return fmt.Errorf("fail to read method negotiation req, err=%w", err)
		}

		authHandler, err := srv.selectAuthMethod(methods)
		if err != nil {
			return err
		}

		if err := protocol.WriteMethodNegotiationReply(authHandler.Method, conn); err != nil {
			return fmt.Errorf("fail to write method negotiation reply, err=%w", err)
		}

		authInfo, err = authHandler.Handle(ctx, conn, srv.authMngr)
		if err != nil {
			return err
		}

		cmd, targetAddr, err := protocol.ReadCommandNegotiationReq(conn, buf)
		if err != nil {
			return fmt.Errorf("fail to read command negotiation req, err=%w", err)
		}

		commandHandler, err := srv.getCommandMethod(cmd)
		if err != nil {
			return err
		}

		to, addr, err := commandHandler.Handle(ctx, authInfo, targetAddr, conn, srv.connMngr)
		if err != nil {
			return err
		}
		if err := protocol.WriteCommandNegotiationReply(conn, buf, addr.String()); err != nil {
			_ = to.Close()
			return err
		}
		return nil
	}

	if err := handshake(); err != nil {
		if _, ok := err.(protocol.NetworkError); !ok {
			_ = conn.Close()
		}
		if socksErr, ok := err.(protocol.SocksError); ok {
			socksErr.SendErrorReply(conn)
		}
		return
	}

}

func (srv *Server) AddAuthHandler(name string) error {
	srv.handlerLock.Lock()
	defer srv.handlerLock.Unlock()

	handler, ok := protocol.AuthHandlers[name]
	if !ok {
		return fmt.Errorf("illeagal auth handler: %s", name)
	}
	srv.authHandlers[handler.Method] = handler
	return nil
}

func (srv *Server) RemoveAuthHandler(name string) error {
	srv.handlerLock.Lock()
	defer srv.handlerLock.Unlock()

	for method, handler := range srv.authHandlers {
		if handler.Name == name {
			delete(srv.authHandlers, method)
			return nil
		}
	}
	return fmt.Errorf("no such auth handler: %s", name)
}

func (srv *Server) AddCommandHandler(name string) error {
	srv.handlerLock.Lock()
	defer srv.handlerLock.Unlock()

	handler, ok := protocol.CommandHandlers[name]
	if !ok {
		return fmt.Errorf("illeagal command handler: %s", name)
	}
	srv.commandHandlers[handler.Method] = handler
	return nil
}

func (srv *Server) RemoveCommandHandler(name string) error {
	srv.handlerLock.Lock()
	defer srv.handlerLock.Unlock()

	for method, handler := range srv.commandHandlers {
		if handler.Name == name {
			delete(srv.commandHandlers, method)
			return nil
		}
	}
	return fmt.Errorf("no such command handler: %s", name)
}

func (srv *Server) selectAuthMethod(methods []byte) (protocol.AuthHandler, error) {
	srv.handlerLock.Lock()
	defer srv.handlerLock.Unlock()

	for m, handler := range srv.authHandlers {
		for _, allowed := range methods {
			if m == allowed {
				return handler, nil
			}
		}
	}

	return protocol.AuthHandler{}, protocol.NoAcceptedMethod
}

func (srv *Server) getCommandMethod(method byte) (protocol.CommandHandler, error) {
	srv.handlerLock.Lock()
	defer srv.handlerLock.Unlock()

	for m, handler := range srv.commandHandlers {
		if m == method {
			return handler, nil
		}
	}

	return protocol.CommandHandler{}, protocol.CommandNotSupported
}

func (srv *Server) Close() {
	_ = srv.listener.Close()
	srv.workers.Range(func(key, value any) bool {
		pipe := value.(*Pipe)
		pipe.cancel()
		return true
	})
}
