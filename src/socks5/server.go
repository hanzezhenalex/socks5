package socks5

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/hanzezhenalex/socks5/src"
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
	config   Config
	connMngr src.ConnectionManager
	authMngr src.AuthManager

	listener net.Listener

	mutex          sync.Mutex
	commanders     map[uint8]Commander
	authenticators map[uint8]Authenticator
}

func NewServer(cfg Config, connMngr src.ConnectionManager, authMngr src.AuthManager, errCh chan error) (*Server, error) {
	srv := &Server{
		config:         cfg,
		connMngr:       connMngr,
		authMngr:       authMngr,
		commanders:     make(map[byte]Commander),
		authenticators: make(map[byte]Authenticator),
	}

	for _, name := range cfg.Auth {
		if err := srv.AddAuthenticator(name); err != nil {
			return nil, err
		}
	}
	for _, name := range cfg.Command {
		if err := srv.AddCommander(name); err != nil {
			return nil, err
		}
	}

	srv.Start(errCh)
	return srv, nil
}

func (srv *Server) Start(errCh chan error) {
	logrus.Infof("start socks server, ip=%s, commands=%s, auth=%s",
		srv.config.Addr(),
		strings.Join(srv.config.Command, ","),
		strings.Join(srv.config.Auth, ","),
	)

	if listener, err := net.Listen("tcp", srv.config.Addr()); err != nil {
		logrus.Errorf("an error happened when starting socks server, err=%s", err.Error())
		errCh <- err
		close(errCh)
		return
	} else {
		srv.listener = listener
	}

	go func() {
		for {
			conn, err := srv.listener.Accept()
			if err != nil {
				if src.ReadOnClosedSocketError(err) {
					err = nil
					logrus.Info("socks server closed")
				} else {
					logrus.Errorf("an error happened when running socks server, err=%s", err.Error())
				}
				errCh <- err
				return
			}
			go srv.onConnection(conn)
		}
	}()
}

func (srv *Server) onConnection(conn net.Conn) {
	ctx := src.NewTraceContext(context.Background())
	tracer := logrus.WithField("id", src.GetIDFromContext(ctx))

	tracer.Infof("new connection from %s", conn.RemoteAddr().String())
	if err := srv.handshake(ctx, conn, tracer); err != nil {
		if _, ok := err.(NetworkError); !ok {
			if socksErr, ok := err.(socksError); ok {
				socksErr.sendErrorReply(conn)
			}
			_ = conn.Close()
		}
		tracer.Errorf("an error happens when handling new connection, err=%s", err.Error())
	}
	tracer.Info("handshake successfully, piping now")
}

func (srv *Server) handshake(ctx context.Context, conn net.Conn, tracer *logrus.Entry) error {
	var buf = make([]byte, maxAddrLen)
	authInfo, err := srv.authenticate(ctx, conn, buf)
	if err != nil {
		return err
	}
	to, err := srv.handleCommand(ctx, conn, buf, authInfo, tracer)
	if err != nil {
		return err
	}
	return srv.connMngr.Pipe(ctx, authInfo, conn, to)
}

func (srv *Server) authenticate(ctx context.Context, conn net.Conn, buf []byte) (src.AuthInfo, error) {
	var (
		authInfo      src.AuthInfo
		authenticator Authenticator
	)
	methods, err := readMethodNegotiationReq(conn, buf)
	if err != nil {
		return authInfo, err
	}

LOOP:
	for _, allowed := range methods {
		for supported, _authenticator := range srv.authenticators {
			if supported == allowed {
				authenticator = _authenticator
				break LOOP
			}
		}
	}

	if authenticator == nil {
		return authInfo, noAcceptedMethod
	}

	if err := writeMethodNegotiationReply(authenticator.Method(), conn, buf); err != nil {
		return authInfo, err
	}
	return authenticator.Handle(ctx, conn)
}

func (srv *Server) handleCommand(
	ctx context.Context,
	conn net.Conn,
	buf []byte,
	authInfo src.AuthInfo,
	tracer *logrus.Entry,
) (net.Conn, error) {
	var commander Commander

	cmd, target, err := readCommandNegotiationReq(conn, buf)
	if err != nil {
		return nil, err
	}

LOOP:
	for supported, _commander := range srv.commanders {
		if supported == cmd {
			commander = _commander
			break LOOP
		}
	}

	if commander == nil {
		return nil, commandNotSupported
	}

	tracer.Infof("socks request: cmd=%s, target=%s", commander.Name(), target.String())
	return commander.Handle(ctx, authInfo, target, conn, buf)
}

func (srv *Server) AddAuthenticator(name string) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	var authenticator Authenticator

	switch name {
	case authNoAuth:
		authenticator = NoAuth{}
	default:
		return fmt.Errorf("illeagal authenticator: %s", name)
	}
	srv.authenticators[authenticator.Method()] = authenticator
	return nil
}

func (srv *Server) RemoveAuthenticator(name string) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	for method, authenticator := range srv.authenticators {
		if authenticator.Name() == name {
			delete(srv.authenticators, method)
			return nil
		}
	}
	return fmt.Errorf("no such authenticator: %s", name)
}

func (srv *Server) AddCommander(name string) error {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()

	var commander Commander

	switch name {
	case cmdConnect:
		commander = NewConnectCommandor(srv.connMngr)
	default:
		return fmt.Errorf("illeagal commander: %s", name)
	}
	srv.commanders[commander.Method()] = commander
	return nil
}

func (srv *Server) Close() {
	_ = srv.listener.Close()
	srv.connMngr.Close()
}
