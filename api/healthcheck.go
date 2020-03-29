package api

import (
	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/gin-gonic/gin"
)

func healthcheck(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/healthcheck", func(c *gin.Context) {
		c.Writer.WriteHeader(204)
	})
}
