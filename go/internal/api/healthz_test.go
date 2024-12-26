package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

func TestHealthz(t *testing.T) {
	log := logging.NewLogger("fatal", false)

	t.Run("db returning valid json from memory succeeds", func(t *testing.T) {
		service := newMockDotaService(getMockQuiz)
		handler := newHealthzHandler(log, service)

		req := httptest.NewRequest("GET", "/healthz", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rr.Code)
		}
	})
	t.Run("bad json fails", func(t *testing.T) {
		service := newMockDotaService(func() map[string]interface{} {
			return map[string]interface{}{
				"this": "shouldn't",
				"work": "at all",
			}
		})

		handler := newHealthzHandler(log, service)

		req := httptest.NewRequest("GET", "/healthz", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// change this after verifying it doesn't work
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %d, with reason: %q", rr.Code, rr.Body.String())
		}
	})

	t.Run("no json fails", func(t *testing.T) {
		service := &MockDotaService{}
		handler := newHealthzHandler(log, service)

		req := httptest.NewRequest("GET", "/healthz", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		// change this after verifying it doesn't work
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code 500, got %d, with reason: %q", rr.Code, rr.Body.String())
		}
	})
}
