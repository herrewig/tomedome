package dota

import (
	"testing"
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
