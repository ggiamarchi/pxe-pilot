package api

import (
	"dev.splitted-desktop.com/horizon/pxe-pilot/model"
	"gopkg.in/gin-gonic/gin.v1"
)

func healthcheck(api *gin.Engine, appConfig *model.AppConfig) {
	api.GET("/healthcheck", func(c *gin.Context) {
		c.Writer.WriteHeader(204)
	})
}
