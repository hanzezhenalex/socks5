package src

import (
	"fmt"

	"github.com/hanzezhenalex/socks5/src/socks5"
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
	config   Config
	socksSrv *socks5.Server
}

func NewAgent(config Config) *Agent {
	return &Agent{
		config: config,
	}
}

func (agent *Agent) Run() error {
	var (
		connMngr ConnectionManager
		authMngr AuthManager
		errCh    = make(chan error)
	)

	switch agent.config.Mode {
	case LocalMode:
		connMngr = NewConnectionManagement()
		authMngr = struct{}{}
	default:
		return fmt.Errorf("%s mode is not supported yet", agent.config.Mode)
	}

	socksSrv, err := socks5.NewServer(agent.config.Socks5Config, connMngr, authMngr, errCh)
	if err != nil {
		return err
	}
	agent.socksSrv = socksSrv

	err = <-errCh
	close(errCh)
	return err
}

func (agent *Agent) Close() {
	agent.socksSrv.Close()
}
