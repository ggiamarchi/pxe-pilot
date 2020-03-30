package api

import (
	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/ggiamarchi/pxe-pilot/service"
	"github.com/gin-gonic/gin"
)

func readConfigurations(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/configurations", func(c *gin.Context) {
		configurations := service.ReadConfigurations(appConfig)
		c.JSON(200, configurations)
	})
}

func showConfiguration(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.GET("/configurations/:name", func(c *gin.Context) {
		name := c.Param("name")
		config, err := service.ReadConfigurationContent(appConfig, name)

		if err != nil {
			switch v := err.(type) {
			case *service.PXEError:
				pxeErrorResponse(c, v)
			default:
				c.Writer.WriteHeader(500)
			}
			return
		}

		c.JSON(200, config)
	})
}

func deployConfiguration(api *gin.RouterGroup, appConfig *model.AppConfig) {
	api.PUT("/configurations/:name/deploy", func(c *gin.Context) {

		var hosts model.HostsQuery

		if err := c.BindJSON(&hosts); err != nil {
			logger.Debug("Invalid data received - %s", err)
			c.Writer.WriteHeader(400)
			return
		}

		resp, err := service.DeployConfiguration(appConfig, c.Param("name"), hosts.Hosts)
		if err != nil {
			switch v := err.(type) {
			case *service.PXEError:
				pxeErrorResponse(c, v)
			default:
				c.Writer.WriteHeader(500)
			}
			return
		}
		c.JSON(200, resp)
	})
}

func pxeErrorResponse(c *gin.Context, err *service.PXEError) {
	code := 500
	switch err.Kind {
	case "NOT_FOUND":
		code = 404
	case "CONFLICT":
		code = 409
	case "BAD_REQUEST":
		code = 400
	default:
		logger.Error("%+v", err)
	}
	c.JSON(code, &struct {
		Message string
	}{
		Message: err.Msg,
	})
}
