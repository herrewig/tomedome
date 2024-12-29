package api

import (
	"net/http"

	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/sirupsen/logrus"
)

// QuizHandler handles the /quiz route. GETs return a random hero quiz JSON
type QuizHandler struct {
	dotes dota.DotaService
	log   *logrus.Entry
}

func newQuizHandler(log *logrus.Entry, dotes dota.DotaService) *QuizHandler {
	return &QuizHandler{
		dotes: dotes,
		log:   log.WithField("handler", "quiz"),
	}
}


func (h *QuizHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Return Quiz JSON from DotaServiceApi
	quiz, err := h.dotes.GetQuizJson()
	if err != nil {
		h.log.WithField("error", err).Fatal("failed to get quiz")
	}

	// Only allow GET requests
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		// Handle CORS
		w.Header().Set("Access-Control-Allow-Origin", getCorsOrigin(r.Header.Get("Origin")))
		w.Write(quiz)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
