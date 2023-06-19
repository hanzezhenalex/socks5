package src

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
)

type Role struct {
	Name string
}

type AuthInfo struct {
	Username string
	Token    string
	Roles    []Role
}

type UserInfo struct {
	Username string
	Password uint32
	Roles    []Role
}

func (info UserInfo) key() string {
	return fmt.Sprintf("%s:%d", info.Username, info.Password)
}

type DataStore interface {
	GetUserInfo(ctx context.Context, username string) (UserInfo, error)
	StoreUserInfo(ctx context.Context, newUser UserInfo) error
}

var (
	userNotExist      = fmt.Errorf("user not exist")
	incorrectPassword = fmt.Errorf("incorrect password")
)

type LocalDataStore struct {
	store sync.Map
}

func (s *LocalDataStore) GetUserInfo(_ context.Context, username string) (UserInfo, error) {
	if info, ok := s.store.Load(username); ok {
		return info.(UserInfo), nil
	}
	return UserInfo{}, userNotExist
}

func (s *LocalDataStore) StoreUserInfo(_ context.Context, newUser UserInfo) error {
	s.store.Store(newUser.Username, newUser)
	return nil
}

type AuthManager interface {
	Login(ctx context.Context, username string, password string) (AuthInfo, error)
	CreateNewUser(ctx context.Context, newUser UserInfo, creator AuthInfo) error
}

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "admin123"
)

type AuthManagement struct {
	store DataStore
}

func (am *AuthManagement) Login(ctx context.Context, username string, password uint32) (AuthInfo, error) {
	user, err := am.store.GetUserInfo(ctx, username)
	if err != nil {
		return AuthInfo{}, err
	}
	if hash(UserInfo{Username: username, Password: password}.key()) != user.Password {
		return AuthInfo{}, incorrectPassword
	}
	authInfo := AuthInfo{
		Username: user.Username,
		Roles:    user.Roles,
	}
	return authInfo, nil
}

func (am *AuthManagement) CreateNewUser(ctx context.Context, newUser UserInfo, creator AuthInfo) error {
	newUser.Password = hash(newUser.key())
	return am.store.StoreUserInfo(ctx, newUser)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}
