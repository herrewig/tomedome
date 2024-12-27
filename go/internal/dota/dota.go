// Service layer for the API. Defines the Client and DotaService interface, and the DotaServiceApi struct.

package dota

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/sirupsen/logrus"
)

// Main service interface for the quiz API
type DotaService interface {
	GetQuizJson() ([]byte, error)
}

// Implementors of this interface load the dota data from a backend
type DotaClient interface {
	GetAllHeroes() ([]Hero, error)
}

// This is is the service api that implements DotaService
type DotaServiceApi struct {
	Client DotaClient
	log    *logrus.Entry
}

type Hero struct {
	Id               int       `json:"id"`
	ShortName        string    `json:"short_name"`
	DisplayName      string    `json:"display_name"`
	Description      string    `json:"description"`
	Lore             string    `json:"lore"`
	AttackType       string    `json:"attack_type"`
	PrimaryAttribute string    `json:"primary_attribute"`
	Abilities        []Ability `json:"abilities"`
}

type Ability struct {
	Id          int    `json:"id"`
	ShortName   string `json:"short_name"`
	DisplayName string `json:"display_name"`
	Button      string `json:"button"`
	Description string `json:"description"`
}

func NewDotaService(log *logrus.Entry, client DotaClient) *DotaServiceApi {
	log.Info("initializing dota service")
	return &DotaServiceApi{
		Client: client,
		log:    log,
	}
}

// Returns the quiz JSON for a random hero
func (s *DotaServiceApi) GetQuizJson() ([]byte, error) {
	s.log.Info("fetching quiz json")

	hero, err := s.GetRandomHero()
	if err != nil {
		s.log.Errorf("Failed to get random hero: %v", err)
		return nil, err
	}
	return hero.GetQuizJson(), nil
}

func (s *DotaServiceApi) GetAllHeroes() ([]Hero, error) {
	s.log.Info("fetching all heroes")

	heroes, err := s.Client.GetAllHeroes()
	if err != nil {
		return nil, err
	}
	return heroes, nil
}

func (s *DotaServiceApi) SerializeDb() ([]byte, error) {
	s.log.Info("serializing db")

	heroes, err := s.Client.GetAllHeroes()
	if err != nil {
		s.log.Errorf("Failed to list heroes: %v", err)
		return nil, err
	}
	j, err := json.Marshal(heroes)
	if err != nil {
		s.log.Errorf("Failed to marshal heroes: %v", err)
		return nil, err
	}
	return j, nil
}

func (h *Hero) PrimaryAttributeString() string {
	switch h.PrimaryAttribute {
	case "agi":
		return "Agility"
	case "str":
		return "Strength"
	case "int":
		return "Intelligence"
	case "all":
		return "Universal"
	default:
		return "Unknown"
	}
}

// Translate hero ability info into a question
func abilityToQuestion(ability Ability) string {
	if ability.Button == "Passive" {
		return "Describe passive ability"
	}
	return fmt.Sprintf("Which ability is mapped to [%s]?", ability.Button)
}

// Get a JSON representation of a hero quiz
func (h *Hero) GetQuizJson() []byte {
	quiz := make(map[string]interface{})
	quiz["displayName"] = h.DisplayName
	quiz["shortName"] = h.ShortName
	quiz["primaryAttribute"] = h.PrimaryAttributeString()
	quiz["attackType"] = h.AttackType
	quiz["description"] = h.Description

	questions := []map[string]string{}

	for _, ability := range h.Abilities {
		questions = append(questions, map[string]string{
			"prompt":      abilityToQuestion(ability),
			"abilityName": ability.DisplayName,
			"answer":      ability.Description,
		})
	}

	quiz["questions"] = questions
	j, err := json.Marshal(quiz)
	if err != nil {
		return nil
	}
	return j
}

func (s *DotaServiceApi) GetRandomHero() (*Hero, error) {
	heroes, err := s.Client.GetAllHeroes()
	if err != nil {
		s.log.Errorf("Failed to list heroes: %v", err)
		return nil, err
	}
	num := len(heroes)
	i := rand.Intn(num)
	return &heroes[i], nil
}
