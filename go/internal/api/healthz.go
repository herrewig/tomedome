package api

import (
	"encoding/json"
	"net/http"

	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/sirupsen/logrus"
)

// HealthzHandler is an http.Handler that returns a 200 OK if the server is healthy
type HealthzHandler struct {
	log   *logrus.Entry
	dotes dota.DotaService
}

func newHealthzHandler(log *logrus.Entry, dotes dota.DotaService) *HealthzHandler {
	return &HealthzHandler{
		log.WithField("handler", "healthz"),
		dotes,
	}
}

// We define the API as healthy if it can read quiz data from the in-memory
// db, unmarshal it, and assert some basic fields are present in the JSON
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Debug("healthz request")

	// Clients need to use GET
	if r.Method != http.MethodGet {
		h.log.WithField("method", r.Method).Info("healthz request with invalid method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// It's our fault if we can't read our own in-memory db
	hero, err := h.dotes.GetQuizJson()
	if err != nil {
		h.log.WithField("error", err).Error("failed to fetch quiz json")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Response
	got := struct {
		DisplayName      string        `json:"displayName"`
		ShortName        string        `json:"shortName"`
		PrimaryAttribute string        `json:"primaryAttribute"`
		AttackType       string        `json:"attackType"`
		Description      string        `json:"description"`
		Questions        []interface{} `json:"questions"`
	}{}

	// It's our fault if our quiz data can't be unmarshaled since we're
	// reading from our own in-memory db that gets created at initialization
	if err := json.Unmarshal([]byte(hero), &got); err != nil {
		h.log.WithField("error", err).Error("failed to unmarshal quiz json")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Make sure all the expected fields are present
	var badJson bool = false
	switch {
	case got.DisplayName == "":
		badJson = true
	case got.ShortName == "":
		badJson = true
	case got.PrimaryAttribute == "":
		badJson = true
	case got.AttackType == "":
		badJson = true
	case got.Description == "":
		badJson = true
	// Quiz should return questions!
	case len(got.Questions) < 1:
		badJson = true
	}

	// It's our fault if the quiz JSON fields are broken
	if badJson {
		h.log.WithField("json", got).Error("quiz json is missing fields")
		http.Error(w, "Internal service error", http.StatusInternalServerError)
	}

	// Respond with a 200 OK if everything is good
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
