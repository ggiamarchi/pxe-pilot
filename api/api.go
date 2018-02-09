package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
	"gopkg.in/gin-gonic/gin.v1"
	yaml "gopkg.in/yaml.v2"
)

// Run runs the PXE Pilot server
func Run(appConfigFile string) {
	logger.Info("Starting PXE Pilot server...")
	appConfig := loadAppConfig(appConfigFile)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler:      api(appConfig),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.ListenAndServe()
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
	rebootHost(api, appConfig)

	return api
}
