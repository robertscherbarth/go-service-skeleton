package main

import (
	"github.com/go-chi/chi"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/robertscherbarth/go-service-skeleton/pkg/app"
	requestLogger "github.com/robertscherbarth/go-service-skeleton/pkg/log"
)

const port = 8080

func main() {
	logger := log.New()
	logger.Formatter = &log.TextFormatter{
		DisableTimestamp: true,
	}

	logger.Infof("starting service ...")
	runningService(logger)

}

func runningService(logger *log.Logger) {
	router := chi.NewRouter()

	//define additional routes
	router.HandleFunc("/", mainHandler)

	app := app.NewApp(&app.Config{HTTPListenPort: port}, router, logger)
	app.CreateRouteConfiguration(requestLogger.NewStructuredLogger(logger))

	app.Start()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `implement me`)
}
