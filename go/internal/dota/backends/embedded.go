// This is a backend for the dota db where the data is baked into the binary at build time.
//
// Data is loaded from the embedded file into memory when the client is created. The data is then
// served from memory.

package backends

import (
	"embed"
	"encoding/json"

	"github.com/herrewig/tomedome/go/internal/assets"
	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/sirupsen/logrus"
)

type EmbeddedDataClient struct {
	log *logrus.Entry
	Db  []dota.Hero
}

func NewEmbeddedDataClient(log *logrus.Entry, filename string) *EmbeddedDataClient {
	newLog := log.WithField("backend", "embedded")
	newLog.Info("loading heroes from embedded data")
	data, _ := getAllHeroDataEmbedded(newLog, filename)

	client := &EmbeddedDataClient{
		log: newLog,
		Db:  data,
	}
	return client
}

func (c *EmbeddedDataClient) GetAllHeroes() ([]dota.Hero, error) {
	c.log.Debug("fetching all heroes from memory")
	return c.Db, nil
}

// Loads hero data from an embedded json file. See https://pkg.go.dev/embed
func getAllHeroDataEmbedded(log *logrus.Entry, fileName string) ([]dota.Hero, error) {
	var assets embed.FS = assets.Assets
	var heroes []dota.Hero

	jsonData, err := assets.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(jsonData), &heroes); err != nil {
		return nil, err
	}
	return heroes, nil
}
