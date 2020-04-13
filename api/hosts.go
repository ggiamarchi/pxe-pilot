package api

import (
	"strconv"

	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/ggiamarchi/pxe-pilot/service"
	"github.com/gin-gonic/gin"
)

func readHosts(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/hosts", func(c *gin.Context) {
		var status bool
		var err error

		statusParam := c.Query("status")
		if statusParam == "" {
			status = true
		} else {
			status, err = strconv.ParseBool(statusParam)
			if err != nil {
				c.Writer.WriteHeader(400)
				return
			}
		}
		hosts := service.ReadHosts(appConfig, status)
		c.JSON(200, hosts)
	})
}

func rebootHost(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.PATCH("/hosts/:name/reboot", func(c *gin.Context) {
		for _, host := range appConfig.Hosts {
			if host.Name == c.Param("name") {
				if service.RebootHost(appConfig, host) != nil {
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

func onHost(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.PATCH("/hosts/:name/on", func(c *gin.Context) {
		for _, host := range appConfig.Hosts {
			if host.Name == c.Param("name") {
				if service.OnHost(appConfig, host) != nil {
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

func offHost(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.PATCH("/hosts/:name/off", func(c *gin.Context) {
		for _, host := range appConfig.Hosts {
			if host.Name == c.Param("name") {
				if service.OffHost(appConfig, host) != nil {
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

func refresh(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.PATCH("/refresh", func(c *gin.Context) {

		if service.Refresh(appConfig) != nil {
			c.Writer.WriteHeader(500)
		}

		c.Writer.WriteHeader(204)
	})
}
