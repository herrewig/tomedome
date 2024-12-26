package backends

import (
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

func TestDotaClientEmbedded(t *testing.T) {
	log := logging.NewLogger("error", false)
	client := NewEmbeddedDataClient(log, "mock_data.json")
	assertFetching(t, client)
}
