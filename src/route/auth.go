package route

import (
	"encoding/json"
	"fmt"
	"github.com/hanzezhenalex/socks5/src/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	tracer = logrus.WithField("component", "auth")
)

const (
	xUserId = "x-user-id"
)

func RegisterAuthManagerEndpoints(router *gin.RouterGroup, auth auth.Manager) {
	router.POST("/user/create", CreateNewUser(auth))
}

type userReq struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

func CreateNewUser(authMngr auth.Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		var req userReq
		if err := json.NewDecoder(context.Request.Body).Decode(&req); err != nil {
			tracer.Errorf("fail to decode CreateNewUser request body, err=%s", err.Error())
			context.Status(http.StatusBadRequest)
			return
		}
		creator, ok := context.Get(xUserId)
		if !ok {
			context.Status(http.StatusInternalServerError)
			return
		}
		newUserInfo := auth.UserInfo{
			Username: req.Username,
			Password: req.Password,
		}
		for _, role := range req.Roles {
			switch role {
			case auth.User:
				newUserInfo.Roles = append(newUserInfo.Roles, auth.RoleUser)
			case auth.Admin:
				newUserInfo.Roles = append(newUserInfo.Roles, auth.RoleAdmin)
			default:
				context.Status(http.StatusBadRequest)
				return
			}
		}

		if err := authMngr.CreateNewUser(context.Request.Context(), newUserInfo, creator.(string)); err != nil {
			if err == auth.NotAuthorize {
				context.Status(http.StatusUnauthorized)
			} else {
				context.Status(http.StatusInternalServerError)
			}
			return
		}
		context.Status(http.StatusOK)
	}
}

func Login(authMngr auth.Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		var req userReq
		if err := json.NewDecoder(context.Request.Body).Decode(&req); err != nil {
			tracer.Errorf("fail to decode Login request body, err=%s", err.Error())
			context.Status(http.StatusBadRequest)
			return
		}
		info, err := authMngr.Login(context.Request.Context(), req.Username, req.Password)
		if err != nil {
			if err == auth.UserNotExist || err == auth.IncorrectPassword {
				context.Status(http.StatusBadRequest)
			} else {
				context.Status(http.StatusInternalServerError)
			}
			return
		}
		context.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", info.Token))
		context.Status(http.StatusOK)
	}
}

func JwtAuth(authMngr auth.Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")
		tokens := strings.Split(authHeader, " ")
		if len(tokens) != 2 {
			context.Status(http.StatusUnauthorized)
			context.Abort()
			return
		}
		jwtToken := tokens[1]
		claims, err := auth.ParseToken(jwtToken)
		if err != nil {
			context.Status(http.StatusUnauthorized)
			context.Abort()
			return
		}
		info, _ := authMngr.GetAuthInfo(context.Request.Context(), claims.Username)
		context.Set(xUserId, info)
		context.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", info.Token))
		context.Next()
	}
}
