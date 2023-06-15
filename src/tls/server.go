package tls

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

type Server struct {
	addr     string
	listener net.Listener
	engine   *gin.Engine
}

func NewServer(addr string) *Server {
	return &Server{
		addr:   addr,
		engine: gin.Default(),
	}
}

func (srv *Server) RouteGroup() *gin.RouterGroup {
	return &srv.engine.RouterGroup
}

func (srv *Server) ListenAndServe(ctx context.Context) error {
	serverCert, caPool, _, err := GenerateCertificates()
	if err != nil {
		return fmt.Errorf("fail to generate certs, %w", err)
	}
	listener, err := tls.Listen("tcp", srv.addr, &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	if err != nil {
		return err
	}

	srv.listener = listener
	ch := make(chan error)

	go func() {
		ch <- srv.engine.RunListener(listener)
	}()

	select {
	case <-ctx.Done():
		_ = srv.listener.Close()
	case err := <-ch:
		return err
	}
	return nil
}
