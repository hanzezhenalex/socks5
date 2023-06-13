package src

import (
	"github.com/gin-gonic/gin"

	"github.com/hanzezhenalex/socks5/src/socks5"
)

type Config struct {
	socksCfg socks5.Config
}

type Agent struct {
	socks      *socks5.Server
	controller *gin.Engine
}

func NewAgent(cfg Config) (*Agent, error) {
	connMngr := &ConnectionManagement{}
	authMngr := &struct{}{}
	socks, err := socks5.NewServer(cfg.socksCfg, connMngr, authMngr)
	if err != nil {
		return nil, err
	}
	return &Agent{
		socks: socks,
	}, nil
}
