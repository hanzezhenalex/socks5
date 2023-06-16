package connection

import (
	"net/http"

	"github.com/hanzezhenalex/socks5/src"

	"github.com/gin-gonic/gin"
)

func RegisterConnectionManagerEndpoints(router *gin.RouterGroup, connMngr Manager) {
	router.GET("/list", ListPipes(connMngr))
}

func ListPipes(connMngr Manager) gin.HandlerFunc {
	return func(context *gin.Context) {
		connections := connMngr.ListConnections(context.Request.Context(), src.AuthInfo{})
		context.JSON(http.StatusOK, connections)
	}
}
