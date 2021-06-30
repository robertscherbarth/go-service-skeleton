package main

import (
	"github.com/robertscherbarth/go-service-skeleton/internal/api"
	"github.com/robertscherbarth/go-service-skeleton/internal/config"
)

// Service as an example of a go micro-service
// @title skeleton-service
// @version 1.0
// @description skeleton service
func main() {
	const configurationFile = "./configs/configuration.yml"

	configuration, err := config.Read("", configurationFile)
	if err != nil {
		panic(err)
	}

	logger, err := config.CreateLogger(configuration.Logger.Level, configuration.Logger.Encoding)

	server := api.NewServer(logger, configuration.HTTP, configuration.Name, configuration.Metrics)

	server.Run()
}
