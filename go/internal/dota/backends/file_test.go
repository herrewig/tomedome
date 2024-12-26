package backends

import (
	"os"
	"testing"

	"github.com/herrewig/tomedome/go/internal/logging"
)

func TestDotaClientFile(t *testing.T) {
	log := logging.NewLogger("error", false)

	// This is set in Makefile
	filePath := os.Getenv("TOMEDOME_DB_FILEPATH")
	if filePath == "" {
		t.Fatal("TOMEDOME_DB_FILEPATH not set")
	}
	client := NewJsonFileClient(log, filePath)
	assertFetching(t, client)
}
