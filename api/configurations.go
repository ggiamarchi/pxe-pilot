package api

import (
	"dev.splitted-desktop.com/horizon/pxe-pilot/logger"
	"dev.splitted-desktop.com/horizon/pxe-pilot/model"
	"dev.splitted-desktop.com/horizon/pxe-pilot/service"
	"gopkg.in/gin-gonic/gin.v1"
)

func readConfigurations(api *gin.Engine, appConfig *model.AppConfig) {
	api.GET("/configurations", func(c *gin.Context) {
		configurations := service.ReadConfigurations(appConfig)
		c.JSON(200, configurations)
	})
}

func deployConfiguration(api *gin.Engine, appConfig *model.AppConfig) {
	api.PUT("/configurations/:name/deploy", func(c *gin.Context) {

		var hosts model.HostsQuery

		if err := c.BindJSON(&hosts); err != nil {
			logger.Debug("Invalid data received - %s", err)
			c.Writer.WriteHeader(400)
			return
		}

		err := service.DeployConfiguration(appConfig, c.Param("name"), hosts.Hosts)
		if err != nil {
			switch v := err.(type) {
			case *service.PXEError:
				pxeErrorResponse(c, v)
			default:
				c.Writer.WriteHeader(500)
			}
			return
		}
		c.Writer.WriteHeader(200)
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
