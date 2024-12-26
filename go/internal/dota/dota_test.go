package dota

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimaryAttributeString(t *testing.T) {
	tests := []struct {
		hero Hero
		want string
	}{
		{Hero{PrimaryAttribute: "str"}, "Strength"},
		{Hero{PrimaryAttribute: "agi"}, "Agility"},
		{Hero{PrimaryAttribute: "int"}, "Intelligence"},
		{Hero{PrimaryAttribute: "all"}, "Universal"},
		{Hero{PrimaryAttribute: "theRealMakaveli"}, "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.hero.PrimaryAttributeString(); got != tt.want {
			t.Errorf("got %q, want %q", got, tt.want)
		}
	}
}

func TestHeroGetQuizJson(t *testing.T) {
	hero := Hero{
		Id:               4,
		ShortName:        "bloodseeker",
		DisplayName:      "Bloodseeker",
		Description:      "Bloodseeker forces difficult decisions on his enemies...",
		Lore:             "Strygwyr the Bloodseeker...",
		AttackType:       "Melee",
		PrimaryAttribute: "agi",
		Abilities: []Ability{
			{
				Id:          5015,
				Button:      "Q",
				ShortName:   "bloodseeker_bloodrage",
				DisplayName: "Bloodrage",
				Description: "Drives bloodseeker...",
			},
			{
				Id:          5017,
				Button:      "Passive",
				ShortName:   "bloodseeker_thirst",
				DisplayName: "Thirst",
				Description: "Bloodseeker is invigorated...",
			},
		},
	}

	got := string(hero.GetQuizJson())
	want := `{"attackType":"Melee","description":"Bloodseeker forces difficult decisions on his enemies...","displayName":"Bloodseeker","primaryAttribute":"Agility","questions":[{"abilityName":"Bloodrage","answer":"Drives bloodseeker...","prompt":"Which ability is mapped to [Q]?"},{"abilityName":"Thirst","answer":"Bloodseeker is invigorated...","prompt":"Describe passive ability"}],"shortName":"bloodseeker"}`

	assert.JSONEq(t, want, got)
}
