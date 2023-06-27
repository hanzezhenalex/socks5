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

	Authorization = "Authorization"
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
			context.String(http.StatusBadRequest, "illegal request body")
			return
		}
		creator, ok := context.Get(xUserId)
		if !ok {
			tracer.Error("no authInfo in context")
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
				context.String(http.StatusBadRequest, "illegal roles of new user: %s", role)
				return
			}
		}

		if err := authMngr.CreateNewUser(context.Request.Context(), newUserInfo, creator.(auth.Info).Username); err != nil {
			if err == auth.NotAuthorize {
				context.String(http.StatusUnauthorized, "permission denied for creator")
			} else if err == auth.IllegalUsernamePassword {
				context.String(http.StatusBadRequest, "illegal username password")
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
			context.String(http.StatusBadRequest, "illegal request body")
			return
		}
		info, err := authMngr.Login(context.Request.Context(), req.Username, req.Password)
		if err != nil {
			if err == auth.UserNotExist || err == auth.IncorrectPassword {
				context.String(http.StatusBadRequest, "incorrect username/password")
			} else {
				context.Status(http.StatusInternalServerError)
			}
			return
		}
		token := fmt.Sprintf("Bearer %s", info.Token)
		context.Writer.Header().Add(Authorization, token)
		context.String(http.StatusOK, token)
	}
}

func JwtAuth(authMngr auth.Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get(Authorization)
		tokens := strings.Split(authHeader, " ")
		if len(tokens) != 2 {
			logrus.Errorf("illegal tokens, header=%s", authHeader)
			context.Status(http.StatusUnauthorized)
			context.Abort()
			return
		}
		jwtToken := tokens[1]
		claims, err := auth.ParseToken(jwtToken)
		if err != nil {
			logrus.Errorf("fail to parse token, illegal tokens, err=%s", err.Error())
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
