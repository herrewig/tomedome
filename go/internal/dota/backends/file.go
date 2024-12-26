// This is a backend for the dota db where the data lives in a plaintext json file on a disk
// filepath.
//
// Data is loaded from the file into memory when the client is created. The data is then served from memory.

package backends

import (
	"encoding/json"
	"io"
	"os"

	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/sirupsen/logrus"
)

type JsonFileClient struct {
	Db  []dota.Hero
	log *logrus.Entry
}

func (c *JsonFileClient) GetAllHeroes() ([]dota.Hero, error) {
	c.log.Debug("fetching all heroes from memory")
	return c.Db, nil
}

func loadHeroesFromFile(log *logrus.Entry, filePath string) ([]dota.Hero, error) {
	log.Info("loading hero data from file")

	// Load heroes from file
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("failed to open file: %v", err)
		return nil, err
	}
	defer file.Close()

	// Read the file contents into a byte slice
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("failed to read file: %v", err)
		return nil, err
	}

	// Unmarshal the JSON into a Go structure
	var heroes []dota.Hero
	if err := json.Unmarshal(bytes, &heroes); err != nil {
		log.Errorf("failed to unmarshal JSON: %v", err)
		return nil, err
	}

	return heroes, nil
}

func NewJsonFileClient(log *logrus.Entry, filePath string) *JsonFileClient {
	newLog := log.WithFields(logrus.Fields{
		"backend":  "json",
		"filePath": filePath,
	})

	db, err := loadHeroesFromFile(newLog, filePath)
	if err != nil {
		log.Fatalf("Failed to load hero data: %v", err)
	}
	return &JsonFileClient{
		Db:  db,
		log: newLog,
	}
}
