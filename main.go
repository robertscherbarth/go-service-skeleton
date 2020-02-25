package main

import (
	"io"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const port = 8080

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Infof("Starting service ...")
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/admin/health", healthCheckHandler)
	http.Handle("/admin/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Infoln("... shutting down")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"path": r.URL.Path, "host": r.URL.Hostname()}).Info()

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `implement me`)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"status": "up"}`)
}
