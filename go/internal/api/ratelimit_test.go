package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

func TestRateLimiting(t *testing.T) {
	log := logging.NewLogger("error", true)
	routes := []string{}

	t.Run("no rate limiting", func(t *testing.T) {
		rateLimited := false
		handler := newHandler(log, rateLimited, routes, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		for i := 0; i < 20; i++ {
			req := newMockRequest("/api/v1/quiz")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status code 200, got %d", rr.Code)
			}
		}
	})

	t.Run("with rate limiting", func(t *testing.T) {
		handler := newLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		for i := 0; i < 20; i++ {
			req := newMockRequest("/api/v1/quiz")
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if i < 10 {
				if rr.Code != http.StatusOK {
					t.Errorf("Request %d: expected status 200, got %d", i, rr.Code)
				}
			}
			if i >= 10 {
				if rr.Code != http.StatusTooManyRequests {
					t.Errorf("Request %d: expected status 429, got %d", i, rr.Code)
				}
			}
		}
	})
}

func TestClientIpMiddleware(t *testing.T) {
	wantIp := "1.2.3.4"

	handler := newLimiterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RemoteAddr != wantIp {
			t.Errorf("Expected %s, got %s", wantIp, r.RemoteAddr)
		}
		w.WriteHeader(http.StatusOK)
	}))
	handler = newClientIpMiddleware(handler)

	for i := 0; i < 20; i++ {
		req := newMockRequest("/api/v1/quiz")
		req.Header.Set("X-Forwarded-For", "172.17.29.38,1.2.3.4,192.168.1.1")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if i < 10 {
			if rr.Code != http.StatusOK {
				t.Errorf("Request %d: expected status 200, got %d", i, rr.Code)
			}
		} else {
			if rr.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d: expected status 429, got %d", i, rr.Code)
			}
		}
	}
}
