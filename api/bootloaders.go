package api

import (
	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/ggiamarchi/pxe-pilot/service"
	"github.com/gin-gonic/gin"
)

func readBootloaders(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/bootloaders", func(c *gin.Context) {
		bootloaders := service.ReadBootloaders(appConfig)
		c.JSON(200, bootloaders)
	})
}
