package api

import (
	"dev.splitted-desktop.com/horizon/pxe-pilot/logger"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/go-pg/pg"
)

// Init initiate declare API endpoints in the framework
func Init(db *pg.DB) *gin.Engine {
	api := gin.New()
	api.Use(logger.APILogger(), gin.Recovery())

	healthcheck(api)

	return api
}
