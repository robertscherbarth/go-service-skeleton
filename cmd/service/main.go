package main

import (
	chiMiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/robertscherbarth/go-service-skeleton/internal/api"
	"github.com/robertscherbarth/go-service-skeleton/internal/config"
	"github.com/robertscherbarth/go-service-skeleton/internal/users"
	"github.com/robertscherbarth/go-service-skeleton/internal/users/adapter"
	"github.com/robertscherbarth/go-service-skeleton/internal/users/ports"
	"go.uber.org/zap"
	"net/http"
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

	initUserDomain(logger, server)

	server.Run()
}

func initUserDomain(logger *zap.Logger, server *api.Server) {
	logger.Info("init user domain with all dependencies")

	userService := users.NewService(logger, adapter.NewInMemoryStore())

	userHttpPort := ports.NewHttp(logger, userService)
	swagger, err := ports.GetSwagger()
	if err != nil {
		logger.Error("-", zap.Error(err))
		return
	}
	swagger.Servers = nil
	validatorMiddleware := chiMiddleware.OapiRequestValidator(swagger)
	ports.HandlerWithOptions(userHttpPort, ports.ChiServerOptions{
		BaseRouter: server.Router,
		Middlewares: []ports.MiddlewareFunc{func(next http.HandlerFunc) http.HandlerFunc {
			return validatorMiddleware(next).ServeHTTP
		}},
	})
}
