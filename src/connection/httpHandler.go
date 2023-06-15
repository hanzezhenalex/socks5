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
		data, err := connMngr.ListConnections(context.Request.Context(), src.AuthInfo{})
		if err != nil {
			context.Status(http.StatusInternalServerError)
		} else {
			context.JSON(http.StatusOK, data)
		}
	}
}
