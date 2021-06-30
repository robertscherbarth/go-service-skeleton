package api

import (
	"context"
	"fmt"
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
	mw "github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/robertscherbarth/go-service-skeleton/pkg/api/middleware"
	"github.com/robertscherbarth/go-service-skeleton/pkg/config"
)

type Server struct {
	logger *zap.Logger

	HTTPListenPort int
	server         *http.Server
	router         chi.Router
}

func NewServer(logger *zap.Logger, config config.HTTP, serviceName string, metrics config.Metrics) *Server {
	router := chi.NewRouter()
	s := &Server{
		logger:         logger,
		HTTPListenPort: config.Port,
		router:         router,
		server: &http.Server{
			Handler: router,
		},
	}

	s.router.Use(
		middleware.RequestLogger(logger, serviceName),
	)

	if metrics.Enabled {
		s.router.Use(middleware.MeasureResponseDuration(metrics.Namespace))
		s.InitializeMetrics(metrics)
	}

	s.InitializeHealth(config.HealthCheck)

	if config.Profiling.Enabled {
		s.InjectProfiling()
	}
	return s
}

func (s *Server) InjectProfiling() {
	s.router.Mount("/debug", mw.Profiler())
}

func (s *Server) InitializeHealth(config config.HealthCheck) {
	s.router.Get(config.Path, s.handleHealthCheck())
	s.logger.Info(fmt.Sprintf("initialize health enpoint to path %s", config.Path))
}

// @Description Returns everytime a 200 response code and a simple json which says {status: up}
// @Tags infrastructure
// @Produce json
// @Success 200
// @Router /health [get]
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		io.WriteString(w, `{"status": "up"}`)
	}
}

func (s *Server) InitializeMetrics(config config.Metrics) {
	s.router.Get(config.Path, s.handleMetrics())
	s.logger.Info(fmt.Sprintf("initialize metrics to path %s", config.Path))
}

// @Description Returns the metrics for prometheus
// @Tags infrastructure
// @Success 200
// @Router /admin/metrics [get]
func (s *Server) handleMetrics() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

func (s *Server) Run() error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.logger.Info("shutting down", zap.String("reason", "received quit signal"))
		s.server.SetKeepAlivesEnabled(false)
		err := s.server.Shutdown(ctx)
		if err != nil {
			s.logger.Panic(err.Error())
		}

		wg.Done()
	}()

	listener, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(s.HTTPListenPort)))
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("server running on port: %d", s.HTTPListenPort))

	err = s.server.Serve(listener)
	if err != http.ErrServerClosed {
		return err
	}

	wg.Wait()
	s.logger.Info("stopped HTTP server")

	return nil
}
