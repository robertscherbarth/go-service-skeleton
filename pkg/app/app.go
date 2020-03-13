package app

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HTTPListenPort int
}

type App struct {
	cfg *Config

	logger *logrus.Logger
	server *http.Server
}

func NewApp(cfg *Config, router *chi.Mux, logger *logrus.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		server: &http.Server{
			Handler: router,
		},
	}
}

func (a *App) Start() error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		a.logger.WithField("reason", "recieved quit signal")
		a.server.SetKeepAlivesEnabled(false)
		err := a.server.Shutdown(ctx)
		if err != nil {
			a.logger.Panic(err.Error())
		}

		wg.Done()
	}()

	listener, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(a.cfg.HTTPListenPort)))
	if err != nil {
		return err
	}
	a.logger.WithField("address", listener.Addr().String()).WithField("port", a.cfg.HTTPListenPort).Info("server listening on address")

	err = a.server.Serve(listener)
	if err != http.ErrServerClosed {
		return err
	}

	wg.Wait()
	a.logger.Info("stopped HTTP server")

	return nil
}

func (a *App) CreateRouteConfiguration(loggerMiddleware func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(loggerMiddleware)
	r.Use(middleware.Recoverer)

	r.Route("/admin", func(r chi.Router) {
		r.HandleFunc("/health", healthCheckHandler)
		r.Handle("/metrics", promhttp.Handler())
	})

	a.server.Handler = r
	return r
}

func (a *App) UpdateRouter(router *chi.Mux) {
	a.server.Handler = router
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, `{"status": "up"}`)
}
