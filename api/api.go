package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ggiamarchi/pxe-pilot/logger"
	"github.com/ggiamarchi/pxe-pilot/model"
	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
)

// Run runs the PXE Pilot server
func Run(appConfigFile string) {
	logger.Info("Starting PXE Pilot server...")
	appConfig := loadAppConfig(appConfigFile)

	logger.Debug("PXE Pilot Configuration loaded : %+v", appConfig)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler:      api(appConfig),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		logger.Error("%s", err)
	}
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

	v1 := api.Group("/v1")

	healthcheck(v1, appConfig)

	readConfigurations(v1, appConfig)
	showConfiguration(v1, appConfig)
	deployConfiguration(v1, appConfig)

	readBootloaders(v1, appConfig)

	readHosts(v1, appConfig)
	rebootHost(v1, appConfig)
	onHost(v1, appConfig)
	offHost(v1, appConfig)

	refresh(v1, appConfig)

	return api
}
