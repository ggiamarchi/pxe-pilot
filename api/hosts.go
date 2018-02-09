package api

import (
	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/ggiamarchi/pxe-pilot/service"
	"gopkg.in/gin-gonic/gin.v1"
)

func readHosts(api *gin.Engine, appConfig *model.AppConfig) {
	api.GET("/hosts", func(c *gin.Context) {
		hosts := service.ReadHosts(appConfig)
		c.JSON(200, hosts)
	})
}

func rebootHost(api *gin.Engine, appConfig *model.AppConfig) {
	api.PATCH("/hosts/:name/reboot", func(c *gin.Context) {
		for _, host := range appConfig.Hosts {
			if host.Name == c.Param("name") {
				if service.RebootHost(host) != nil {
					c.Writer.WriteHeader(409)
					return
				}
				c.Writer.WriteHeader(204)
				return
			}
		}
		c.Writer.WriteHeader(404)
	})
}
