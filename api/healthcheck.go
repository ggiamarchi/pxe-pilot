package api

import (
	"github.com/ggiamarchi/pxe-pilot/model"
	"gopkg.in/gin-gonic/gin.v1"
)

func healthcheck(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/healthcheck", func(c *gin.Context) {
		c.Writer.WriteHeader(204)
	})
}
