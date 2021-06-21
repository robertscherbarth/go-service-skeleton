package api

import (
	"context"
	"fmt"
	"github.com/robertscherbarth/go-service-skeleton/pkg/config"
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
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger

	HTTPListenPort int
	server         *http.Server
	router         chi.Router
}

func NewServer(logger *zap.Logger, config config.HTTP) *Server {
	router := chi.NewRouter()
	s := &Server{
		logger:         logger,
		HTTPListenPort: config.Port,
		router:         router,
		server: &http.Server{
			Handler: router,
		},
	}

	if config.Profiling.Enabled {
		s.InjectProfiling()
	}
	return s
}

func (s *Server) InjectProfiling() {
	s.router.Mount("/debug", mw.Profiler())
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
