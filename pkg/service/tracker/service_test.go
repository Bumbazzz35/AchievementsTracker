package tracker

import (
	"testing"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
)

type mockLocalizationRepository struct{}

func (m *mockLocalizationRepository) GetLocalizedAdvancement(id string) (advancement.LocalizedAdvancement, bool) {
	knownMap := map[string]bool{
		"minecraft:adventure/adventuring_time": true,
		"minecraft:story/root":                 true,
		"minecraft:story/mine_stone":           true,
		"minecraft:story/smelt_iron":           true,
		"minecraft:story/upgrade_tools":        true,
	}

	if knownMap[id] {
		return advancement.LocalizedAdvancement{}, true
	}
	return advancement.LocalizedAdvancement{}, false
}

func (m *mockLocalizationRepository) GetAllAdvancementIDs() []string {
	return []string{
		"minecraft:adventure/adventuring_time",
		"minecraft:story/root",
		"minecraft:story/mine_stone",
		"minecraft:story/smelt_iron",
		"minecraft:story/upgrade_tools",
	}
}

func TestGetCurrentAdvancements(t *testing.T) {
	locRepo := &mockLocalizationRepository{}

	tests := []struct {
		name     string
		data     map[string]advancement.Advancement
		expected int
	}{
		{"empty map", map[string]advancement.Advancement{}, 0},
		{"all achievements false", map[string]advancement.Advancement{
			"minecraft:adventure/adventuring_time": {Done: false},
			"minecraft:story/mine_stone":           {Done: false},
			"minecraft:story/upgrade_tools":        {Done: false},
		}, 0},
		{"only recipes", map[string]advancement.Advancement{
			"minecraft:recipes/redstone/lever":       {Done: true},
			"minecraft:recipes/combat/iron_leggings": {Done: true},
			"minecraft:recipes/tools/stone_axe":      {Done: false},
		}, 0},
		{"only real achviements all done", map[string]advancement.Advancement{
			"minecraft:adventure/adventuring_time": {Done: true},
			"minecraft:story/root":                 {Done: true},
			"minecraft:story/mine_stone":           {Done: true},
			"minecraft:story/smelt_iron":           {Done: true},
			"minecraft:story/upgrade_tools":        {Done: true},
		}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCurrentAdvancements(tt.data, locRepo)
			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}
