package agent

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/hanzezhenalex/socks5/src"
	"github.com/hanzezhenalex/socks5/src/connection"
	"github.com/hanzezhenalex/socks5/src/socks5"
	tlsUtil "github.com/hanzezhenalex/socks5/src/tls"
)

type Mode string

const (
	LocalMode   = "local"
	ClusterMode = "cluster"
)

type Config struct {
	Mode              Mode
	ControlServerPort string
	Socks5Config      socks5.Config
}

type Agent struct {
	config        Config
	socksSrv      *socks5.Server
	controlServer *tlsUtil.Server
	closed        atomic.Bool
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
	}
}

func (agent *Agent) Run() error {
	var (
		connMngr        connection.Manager
		authMngr        src.AuthManager
		socksErrCh      = make(chan error, 1)
		controlSrvErrCh = make(chan error, 1)
		wg              sync.WaitGroup
		runningErr      error
	)

	switch agent.config.Mode {
	case LocalMode:
		connMngr = connection.NewConnectionManagement()
		authMngr = struct{}{}
	default:
		return fmt.Errorf("%s mode is not supported yet", agent.config.Mode)
	}

	socksSrv, err := socks5.NewServer(agent.config.Socks5Config, connMngr, authMngr)
	if err != nil {
		return err
	}
	agent.socksSrv = socksSrv
	go func() {
		socksErrCh <- agent.socksSrv.Start()
		close(socksErrCh)
		wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	agent.controlServer = tlsUtil.NewServer(
		fmt.Sprintf("%s:%s", agent.config.Socks5Config.IP, agent.config.ControlServerPort))
	go func() {
		agent.startControlServer(ctx, connMngr, controlSrvErrCh)
		close(controlSrvErrCh)
		wg.Done()
	}()

	wg.Add(2)

	select {
	case runningErr = <-socksErrCh:
	case runningErr = <-controlSrvErrCh:
	}

	agent.Close()
	cancel()
	wg.Wait()
	return runningErr
}

func (agent *Agent) startControlServer(ctx context.Context, connMngr connection.Manager, errCh chan error) {
	routeGroup := agent.controlServer.RouteGroup()
	v1 := routeGroup.Group("/v1")
	connection.RegisterConnectionManagerEndpoints(v1.Group("/connection"), connMngr)

	errCh <- agent.controlServer.ListenAndServe(ctx)
}

func (agent *Agent) Close() {
	if agent.closed.CompareAndSwap(false, true) {
		agent.socksSrv.Close()
	}
}
