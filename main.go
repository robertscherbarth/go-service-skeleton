package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/robertscherbarth/go-service-skeleton/pkg/app"
)

type ServiceConfig struct {
	Port      int    `envconfig:"port" default:"8080"`
	LogFormat string `default:"text"`
}

func main() {
	var serviceConfig ServiceConfig
	err := envconfig.Process("service", &serviceConfig)
	if err != nil {
		fmt.Printf("can't map env to config with err: %v\n", err)
		os.Exit(1)
	}

	logger := createLogger(serviceConfig.LogFormat)
	logger.Infof("starting service ...")
	runService(logger, serviceConfig)

}

func runService(logger *log.Logger, config ServiceConfig) {
	app := app.NewApp(&app.Config{HTTPListenPort: config.Port}, logger)
	app.CreateRouteConfiguration()

	//TODO: define additional routes
	app.AddRoute("/", mainHandler)

	app.Start()
}

func createLogger(logFormat string) *log.Logger {
	logger := log.New()
	if logFormat == "text" {
		logger.Formatter = &log.TextFormatter{
			DisableTimestamp: true,
		}
		return logger
	}
	logger.Formatter = &log.JSONFormatter{
		DisableTimestamp: true,
	}
	return logger
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `implement me`)
}
