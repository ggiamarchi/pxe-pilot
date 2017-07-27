package api

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func healthcheck(api *gin.Engine) {
	api.GET("/healthcheck", func(c *gin.Context) {
		c.Writer.WriteHeader(204)
	})
}
