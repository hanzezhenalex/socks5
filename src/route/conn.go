package route

import (
	"github.com/hanzezhenalex/socks5/src/auth"
	"net/http"

	"github.com/hanzezhenalex/socks5/src/connection"

	"github.com/gin-gonic/gin"
)

func RegisterConnectionManagerEndpoints(router *gin.RouterGroup, connMngr connection.Manager, auth auth.Manager) {
	router.GET("/list", ListPipes(connMngr, auth))
}

func ListPipes(connMngr connection.Manager, _ auth.Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		_, ok := context.Get(xUserId)
		if !ok {
			context.Status(http.StatusInternalServerError)
			return
		}

		connections := connMngr.ListConnections(context.Request.Context(), auth.Info{})
		context.JSON(http.StatusOK, connections)
	}
}
