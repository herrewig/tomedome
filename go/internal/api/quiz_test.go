package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

func TestQuizHappyPath(t *testing.T) {
	service := newMockDotaService(getMockQuiz)
	handler := newQuizHandler(logging.NewLogger("error", false), service)

	req := newMockRequest("/api/v1/quiz")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	want := getMockQuiz()
	got := rr.Body.String()
	if reflect.DeepEqual(want, got) {
		t.Errorf("Expected body %q, got %q", want, got)
	}
}
