package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

func newMockDotaService(quizFunc func() map[string]interface{}) *MockDotaService {
	return &MockDotaService{quizFunc}
}

type MockDotaService struct {
	quizFunc func() map[string]interface{}
}

func (m *MockDotaService) GetQuizJson() ([]byte, error) {
	if m.quizFunc == nil {
		return nil, errors.New("no quizFunc defined")
	}

	j, err := json.Marshal(m.quizFunc())
	if err != nil {
		return nil, err
	}
	return j, nil
}

func getMockQuiz() map[string]interface{} {
	return map[string]interface{}{
		"displayName":      "mock",
		"shortName":        "yeah",
		"primaryAttribute": "ing",
		"attackType":       "yeah",
		"description":      "bird",
		"questions": []map[string]string{
			{
				"prompt":      "yeah",
				"abilityName": "yeahhh",
				"answer":      "yeahhhhh",
			},
		},
	}
}

func TestServerControls(t *testing.T) {
	log := logging.NewLogger("error", false)
	t.Run("params not allowed", func(t *testing.T) {
		middleware := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		handler := newParamValidationMiddleware(log, middleware)

		req := httptest.NewRequest("GET", "/healthz?foo=bar", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %d", rr.Code)
		}
	})

	t.Run("CORS origin", func(t *testing.T) {
		tests := []struct {
			url  string
			want string
		}{
			// Legit localdev
			{
				url:  "http://localhost:8080",
				want: "http://localhost:8080",
			},
			// Local dev Different port
			{
				url:  "http://localhost:8081",
				want: "http://localhost:8080",
			},
			// Local dev no schema
			{
				url:  "localhost:8080",
				want: "http://localhost:8080",
			},
			// Prod URL - this is the one exception to returning http://localhost:8080
			{
				url:  "https://api.tomedome.io",
				want: "https://api.tomedome.io",
			},
			// Prod URL different port
			{
				url:  "https://tomedome.herrewig.dev:4242",
				want: "http://localhost:8080",
			},
			// Prod URL no schema
			{
				url:  "tomedome.herrewig.dev",
				want: "http://localhost:8080",
			},
			// Prod URL no TLS
			{
				url:  "http://tomedome.herrewig.dev",
				want: "http://localhost:8080",
			},
			// Random URL
			{
				url:  "http://lolpwnd.net",
				want: "http://localhost:8080",
			},
			// Empty string
			{
				url:  "",
				want: "http://localhost:8080",
			},
		}
		for _, test := range tests {
			if got := getCorsOrigin(test.url); got != test.want {
				t.Errorf("%q got %q, want %q", test.url, got, test.want)
			}
		}
	})
}
