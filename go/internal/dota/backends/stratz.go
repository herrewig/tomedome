// This is a backend for the dota db where the data is fetched from the public Stratz API.
//
// When the StratzClient is created, https://api.stratz.com/graphiql is called and fetches all data.
// Data is then served from memory.


package backends

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/herrewig/tomedome/go/internal/dota"
	"github.com/sirupsen/logrus"
)

type StratzClient struct {
	apiKey  string
	baseURL string
	log     *logrus.Entry
	db      []dota.Hero
}

// Initializes stratz client as well as loads full dataset from
// the Stratz API
func NewStratzClient(log *logrus.Entry, apiKey string) *StratzClient {
	const url = "https://api.stratz.com/graphql"
	var err error

	newLog := log.WithFields(logrus.Fields{
		"backend":   "stratz",
		"stratzUrl": url,
	})
	if apiKey == "" {
		newLog.Fatal("API key is required")
	}

	client := &StratzClient{
		apiKey:  apiKey,
		baseURL: url,
		log:     newLog,
	}
	client.db, err = client.fetchAllHeroes()
	if err != nil {
		newLog.Fatalf("Failed to load hero data: %v", err)
	}
	return client
}

type AbilityResponse struct {
	Id      int `json:"abilityId"`
	Slot    int `json:"slot"`
	Ability struct {
		Name     string `json:"name"`
		Language struct {
			DisplayName string   `json:"displayName"`
			Description []string `json:"description"`
		} `json:"language"`
		Stat struct {
			Behavior int       `json:"behavior"`
			ManaCost []float64 `json:"manaCost"`
			MaxLevel int       `json:"maxLevel"`
		} `json:"stat"`
	} `json:"ability"`
}

type HeroResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string
	ShortName   string `json:"shortName"`
	Stats       struct {
		AttackType       string `json:"attackType"`
		PrimaryAttribute string `json:"primaryAttribute"`
	} `json:"stats"`
	Language struct {
		DisplayName string `json:"displayName"`
		Lore        string `json:"lore"`
		Hype        string `json:"hype"`
	} `json:"language"`
	Abilities []AbilityResponse `json:"abilities"`
}

type AllHeroesResponse struct {
	Data struct {
		Constants struct {
			Heroes []HeroResponse `json:"heroes"`
		} `json:"constants"`
	} `json:"data"`
}

func (c *StratzClient) Query(query string) (string, error) {
	c.log.Info("querying stratz API")

	// Create a new POST request
	q := `{"query": "` + query + `"}`
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer([]byte(q)))
	if err != nil {
		c.log.WithField("error", err).Error("error creating request")
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "STRATZ_API")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.log.WithField("error", err).Error("error performing request")
		return "", err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.WithField("error", err).Error("error reading response body")
		return "", err
	}

	// Check if the request was successful
	if resp.StatusCode == http.StatusOK {
		c.log.Debug("request successful")
	} else {
		return "", fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return string(body), nil
}

func PrepareQuery(query string) string {
	newStr := strings.ReplaceAll(query, "\n", "")
	newStr = strings.ReplaceAll(newStr, "\r", "")
	newStr = strings.ReplaceAll(newStr, "\t", " ")
	return newStr
}

func (c *StratzClient) Load() error {
	return nil
}

func (c *StratzClient) getAllHeroData() (*AllHeroesResponse, error) {
	c.log.Debug("preparing query")

	query := PrepareQuery(`{
		constants {
			heroes {
				id
				name
				shortName
				stats {
					attackType
					primaryAttribute
				}
				language {
					displayName
					lore
					hype
				}
				abilities {
					abilityId
					slot
					ability {
						name
						stat {
							behavior
							castPoint
							castRange
							manaCost
							maxLevel
						}
						language {
							displayName
							description
						}
					}
				}
			}
		}
	}`)

	c.log.Debug("querying stratz API")
	body, err := c.Query(query)
	if err != nil {
		return nil, err
	}

	c.log.Debug("unmarshalling API response")
	respJson := AllHeroesResponse{}
	err = json.Unmarshal([]byte(body), &respJson)
	if err != nil {
		c.log.Errorf("Error unmarshalling response: %v", err)
		return nil, err
	}

	return &respJson, nil
}

func (c *StratzClient) fetchAllHeroes() ([]dota.Hero, error) {
	// Translate the response to a list of Hero structs
	respJson, err := c.getAllHeroData()
	if err != nil {
		c.log.Errorf("error getting hero data: %v", err)
		return nil, err
	}

	heroes := []dota.Hero{}
	for _, h := range respJson.Data.Constants.Heroes {
		heroes = append(heroes, responseToHero(h))
	}
	return heroes, nil
}

func (c *StratzClient) GetAllHeroes() ([]dota.Hero, error) {
	c.log.Debug("fetching all heroes from memory")
	return c.db, nil
}

/*
Behavior Name                      Value   Description
-----------------------------------------------------------
DOTA_ABILITY_BEHAVIOR_NONE          0       No specific behavior.
DOTA_ABILITY_BEHAVIOR_HIDDEN        1       Hidden ability; not visible to the user.
DOTA_ABILITY_BEHAVIOR_PASSIVE       2       Passive ability; no activation required.
DOTA_ABILITY_BEHAVIOR_NO_TARGET     4       Ability doesn't target anything (e.g., Phase Shift).
DOTA_ABILITY_BEHAVIOR_UNIT_TARGET   8       Ability targets a unit.
DOTA_ABILITY_BEHAVIOR_POINT         16      Ability targets a location.
DOTA_ABILITY_BEHAVIOR_AOE           32      Displays an area of effect (AoE) indicator.
DOTA_ABILITY_BEHAVIOR_CHANNELLED    64      Ability is channeled.
DOTA_ABILITY_BEHAV
*/
func slotToButton(ability AbilityResponse) string {
	// Check if it's a passive ability (or something else weird or exceptional)
	if ability.Ability.Stat.ManaCost == nil &&
		ability.Ability.Stat.MaxLevel == 1 {
		return "Passive"
	}

	switch ability.Ability.Stat.Behavior {
	case 0:
		return "Does nothing"
	case 1:
		return "Hidden"
	case 2, 66:
		return "Passive"
	}

	// Otherwise, use the slot number to determine the button
	switch ability.Slot {
	case 1:
		return "Q"
	case 2:
		return "W"
	case 3:
		return "E"
	case 4:
		return "F"
	case 5:
		return "D"
	case 6:
		return "R"
	case 7:
		return "T"
	case 8:
		return "G"
	default:
		return "Unknown"
	}
}

func responseToHero(r HeroResponse) dota.Hero {
	abilities := []dota.Ability{}

	for _, resp := range r.Abilities {
		// Skip abilities with no description
		if len(resp.Ability.Language.Description) < 1 {
			continue
		}

		abilities = append(abilities, dota.Ability{
			Id:          resp.Id,
			ShortName:   resp.Ability.Name,
			DisplayName: resp.Ability.Language.DisplayName,
			Description: strings.Join(resp.Ability.Language.Description, "\n"),
			Button:      slotToButton(resp),
		})
	}

	return dota.Hero{
		Id:               r.Id,
		ShortName:        r.ShortName,
		DisplayName:      r.Language.DisplayName,
		Description:      r.Language.Hype,
		AttackType:       r.Stats.AttackType,
		PrimaryAttribute: r.Stats.PrimaryAttribute,
		Lore:             r.Language.Lore,
		Abilities:        abilities,
	}
}
