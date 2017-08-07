package api

import (
	"fmt"
	"io/ioutil"

	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
	"gopkg.in/gin-gonic/gin.v1"
	yaml "gopkg.in/yaml.v2"
)

// Run runs the PXE Pilot server
func Run(appConfigFile string) {
	logger.Info("Starting PXE Pilot server...")
	appConfig := loadAppConfig(appConfigFile)
	api(appConfig).Run(fmt.Sprintf(":%d", appConfig.Server.Port))
}

func loadAppConfig(file string) *model.AppConfig {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	c := model.AppConfig{}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		panic(err)
	}

	return &c
}

func api(appConfig *model.AppConfig) *gin.Engine {
	api := gin.New()
	api.Use(logger.APILogger(), gin.Recovery())

	healthcheck(api, appConfig)

	readConfigurations(api, appConfig)
	deployConfiguration(api, appConfig)

	readHosts(api, appConfig)

	return api
}
