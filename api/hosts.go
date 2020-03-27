package api

import (
	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/ggiamarchi/pxe-pilot/service"
	"github.com/gin-gonic/gin"
)

func readHosts(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/hosts", func(c *gin.Context) {
		hosts := service.ReadHosts(appConfig, true)
		c.JSON(200, hosts)
	})
}

func rebootHost(api *gin.RouterGroup, appConfig *model.AppConfig) {
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

func discovery(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.PATCH("/discovery", func(c *gin.Context) {

		if service.Discovery(appConfig) != nil {
			c.Writer.WriteHeader(500)
		}

		c.Writer.WriteHeader(204)
	})
}
