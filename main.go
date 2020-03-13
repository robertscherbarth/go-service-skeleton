package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	requestLogger "github.com/robertscherbarth/go-service-skeleton/pkg/log"
)

const port = 8080

func main() {
	r := chi.NewRouter()

	logger := log.New()
	logger.Formatter = &log.TextFormatter{
		DisableTimestamp: true,
	}

	r.Use(middleware.RealIP)
	r.Use(requestLogger.NewStructuredLogger(logger))
	r.Use(middleware.Recoverer)

	logger.Infof("Starting service ...")
	r.HandleFunc("/", mainHandler)

	r.Route("/admin", func(r chi.Router) {
		r.HandleFunc("/health", healthCheckHandler)
		r.Handle("/metrics", promhttp.Handler())
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(port), r); err != nil {
		logger.Infoln("... shutting down")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `implement me`)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"status": "up"}`)
}
