package localization

import (
	_ "embed"
	"encoding/json"

	"github.com/Bumbazzz35/AchievementTracker/internal/domain/advancement"
)

//go:embed lang_ru.json
var langData []byte

type localizationRepository struct {
	entries map[string]entry
}

func (repo *localizationRepository) GetLocalizedAdvancement(id string) (advancement.LocalizedAdvancement, bool) {
	e, ok := repo.entries[id]
	if !ok {
		return advancement.LocalizedAdvancement{}, false
	}

	return advancement.LocalizedAdvancement{ID: id, Title: e.Title, Description: e.Description, Icon: e.Icon, Category: e.Category}, true
}

func (repo *localizationRepository) GetAllAdvancementIDs() []string {
	ids := make([]string, 0, len(repo.entries))
	for id := range repo.entries {
		ids = append(ids, id)
	}
	return ids
}

func NewRepository() (*localizationRepository, error) {
	repo := localizationRepository{}

	err := json.Unmarshal(langData, &repo.entries)
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

type entry struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
}
