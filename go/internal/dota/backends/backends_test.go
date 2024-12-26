package backends

import (
	"fmt"
	"log"
	"testing"

	"github.com/herrewig/tomedome/go/internal/dota"
)

// Helper test function to assert the different implmentations of DotaClient
// have the same behavior.
func assertFetching(t *testing.T, client dota.DotaClient) {
	t.Run("GetAllHeroes", func(t *testing.T) {
		heroes, err := client.GetAllHeroes()
		if err != nil {
			t.Errorf("Failed to list heroes: %v", err)
		}
		if len(heroes) == 0 {
			t.Error("Expected heroes to be non-empty")
		}
		t.Run("has descriptions", func(t *testing.T) {
			for _, hero := range heroes {
				if hero.Description == "" {
					t.Errorf("Expected hero %d to have a description", hero.Id)
				}
			}
		})
		t.Run("has abilities", func(t *testing.T) {
			for _, hero := range heroes {
				if len(hero.Abilities) == 0 {
					t.Errorf("Expected hero %d to have abilities", hero.Id)
				}
			}
		})
		t.Run("has primary attribute set", func(t *testing.T) {
			for _, hero := range heroes {
				if hero.PrimaryAttributeString() == "" {
					t.Error("Expected hero to have primary attribute")
				}
			}
		})
		t.Run("has primary attribute set", func(t *testing.T) {
			for _, hero := range heroes {
				if hero.AttackType == "" {
					t.Error("Expected hero to have attack type")
					fmt.Println(hero)
					log.Fatal("Expected hero to have attack type")
				}
			}
		})
	})
}
