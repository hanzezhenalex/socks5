package socks5

import (
	"context"
	"github.com/hanzezhenalex/socks5/src/auth"
	"net"
)

const (
	noAuth     = uint8(0x00)
	userPasswd = uint8(0x01)

	authNoAuth     = "noAuth"
	authUserPasswd = "usernamePassword"
)

type Authenticator interface {
	Name() string
	Method() uint8
	Handle(ctx context.Context, conn net.Conn) (auth.Info, error)
}

var emptyAuthInfo = auth.Info{}

type NoAuth struct{}

func (n NoAuth) Name() string {
	return authNoAuth
}

func (n NoAuth) Method() uint8 {
	return noAuth
}

func (n NoAuth) Handle(_ context.Context, _ net.Conn) (auth.Info, error) {
	return emptyAuthInfo, nil
}

type UsernamePassword struct {
	auth auth.Manager
}

func (up UsernamePassword) Name() string {
	return authUserPasswd
}

func (up UsernamePassword) Method() uint8 {
	return userPasswd
}

func (up UsernamePassword) Handle(ctx context.Context, conn net.Conn) (auth.Info, error) {
	username, password, err := readUsernamePasswordAuthRequest(conn)
	if err != nil {
		return auth.Info{}, err
	}
	info, err := up.auth.Login(ctx, username, password)
	if err != nil {
		if err == auth.UserNotExist || err == auth.IncorrectPassword {
			return auth.Info{}, incorrectUsernamePassword
		} else {
			return auth.Info{}, internalError
		}
	}
	if err := writeUsernamePasswordReply(conn); err != nil {
		return auth.Info{}, err
	}
	return info, nil
}
