package api

import (
	"dev.splitted-desktop.com/horizon/pxe-pilot/model"
	"dev.splitted-desktop.com/horizon/pxe-pilot/service"
	"gopkg.in/gin-gonic/gin.v1"
)

func readHosts(api *gin.Engine, appConfig *model.AppConfig) {
	api.GET("/hosts", func(c *gin.Context) {
		hosts := service.ReadHosts(appConfig)
		c.JSON(200, hosts)
		c.Writer.WriteHeader(200)
	})
}
