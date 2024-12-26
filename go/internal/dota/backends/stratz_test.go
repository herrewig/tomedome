package backends

import (
	"os"
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

// Calls the actual Stratz GraphQL API, so we only run
// on integration test runs
func TestDotaClientStratz_Integration(t *testing.T) {
	log := logging.NewLogger("error", false)

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	apiKey := os.Getenv("TOMEDOME_STRATZ_API_KEY")
	if apiKey == "" {
		t.Fatal("TOMEDOME_STRATZ_API_KEY not set")
	}

	client := NewStratzClient(log, apiKey)
	assertFetching(t, client)
}
