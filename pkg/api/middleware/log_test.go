package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	t.Run("should have have access log fields in output", func(t *testing.T) {
		core, logs := observer.New(zap.InfoLevel)
		logger := zap.New(core)

		requestLoggerHandler := RequestLogger(logger, "test-service")
		r := chi.NewRouter()
		r.Use(requestLoggerHandler)
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(r)
		defer ts.Close()

		_, err := http.Get(ts.URL)
		if err != nil {
			t.Error("can't reach http test server")
		}

		if logs.Len() != 1 {
			t.Errorf("no logs are available or to many logs are printed")
		}

		entry := logs.All()[0]
		for _, val := range entry.Context {
			if val.Key == "log_type" {
				if val.String != "access" {
					t.Errorf("expected log_type to be access got %v", val.String)
				}
			}
			if val.Key == "application_type" {
				if val.String != "service" {
					t.Errorf("expected application_type to be service got %v", val.String)
				}
			}
			if val.Key == "service" {
				if val.String != "test-service" {
					t.Errorf("expected service to be test-service got %v", val.String)
				}
			}
		}
	})
}
