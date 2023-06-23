package auth

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"hash/fnv"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Role struct {
	Name string
}

const (
	Admin = "admin"
	User  = "user"
)

var (
	RoleAdmin = Role{Name: Admin}
	RoleUser  = Role{Name: User}
)

type Info struct {
	Username string
	Token    string
	Roles    []Role
}

type UserInfo struct {
	Username string
	Password string
	Roles    []Role
}

func (info UserInfo) IsAdmin() bool {
	for _, role := range info.Roles {
		if role == RoleAdmin {
			return true
		}
	}
	return false
}

func (info UserInfo) key() string {
	return fmt.Sprintf("%s:%s", info.Username, info.Password)
}

type DataStore interface {
	GetUserInfo(ctx context.Context, username string) (UserInfo, error)
	StoreUserInfo(ctx context.Context, newUser UserInfo) error
}

var (
	UserNotExist      = fmt.Errorf("user not exist")
	IncorrectPassword = fmt.Errorf("incorrect password")
	NotAuthorize      = fmt.Errorf("not authorized")

	TokenExpired = fmt.Errorf("token expired")
	TokenInvalid = fmt.Errorf("invalid token")
)

type LocalDataStore struct {
	store sync.Map
}

func (s *LocalDataStore) GetUserInfo(_ context.Context, username string) (UserInfo, error) {
	if info, ok := s.store.Load(username); ok {
		return info.(UserInfo), nil
	}
	return UserInfo{}, UserNotExist
}

func (s *LocalDataStore) StoreUserInfo(_ context.Context, newUser UserInfo) error {
	s.store.Store(newUser.Username, newUser)
	return nil
}

type Manager interface {
	Login(ctx context.Context, username string, password string) (Info, error)
	CreateNewUser(ctx context.Context, newUser UserInfo, creator string) error
	GetAuthInfo(ctx context.Context, username string) (Info, error)
}

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
	defaultSecretKey     = "secret-key"
)

type LocalManagement struct {
	secretKey string
	store     DataStore
}

func NewLocalManagement() *LocalManagement {
	mngr := &LocalManagement{
		secretKey: defaultSecretKey,
		store:     &LocalDataStore{},
	}
	_ = mngr.store.StoreUserInfo(context.Background(), UserInfo{
		Username: DefaultAdminUsername,
		Password: hash(DefaultAdminUsername, DefaultAdminPassword),
		Roles:    []Role{RoleAdmin, RoleUser},
	})
	return mngr
}

func (am *LocalManagement) GetAuthInfo(ctx context.Context, username string) (Info, error) {
	return am.getAuthInfo(ctx, username, nil)
}

func (am *LocalManagement) getAuthInfo(ctx context.Context, username string, verify func(password string) bool) (Info, error) {
	user, err := am.store.GetUserInfo(ctx, username)
	if err != nil {
		return Info{}, err
	}

	if verify != nil && !verify(user.Password) {
		return Info{}, IncorrectPassword
	}

	token, err := MakeToken(username)
	if err != nil {
		return Info{}, err
	}

	authInfo := Info{
		Username: user.Username,
		Roles:    user.Roles,
		Token:    token,
	}
	return authInfo, nil
}

func (am *LocalManagement) Login(ctx context.Context, username string, password string) (Info, error) {
	return am.getAuthInfo(ctx, username, func(target string) bool {
		return hash(username, password) == target
	})
}

func (am *LocalManagement) CreateNewUser(ctx context.Context, newUser UserInfo, creator string) error {
	creatorInfo, err := am.store.GetUserInfo(ctx, creator)
	if err != nil {
		return fmt.Errorf("fail to get creator info, %w", err)
	}
	if !creatorInfo.IsAdmin() {
		return NotAuthorize
	}
	newUser.Password = hash(newUser.Username, newUser.Password)
	return am.store.StoreUserInfo(ctx, newUser)
}

func hash(username string, password string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(fmt.Sprintf("%s:%s", username, password)))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}

type JwtClaims struct {
	Username string
	jwt.RegisteredClaims
}

func MakeToken(username string) (tokenString string, err error) {
	claim := JwtClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // HS256
	return token.SignedString([]byte(defaultSecretKey))
}

func ParseToken(s string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(s, &JwtClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(defaultSecretKey), nil
	})
	if err != nil {
		logrus.Errorf("fail to parse token, err=%s", err.Error())
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			}
			return nil, TokenInvalid
		}
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
