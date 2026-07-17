package localization

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
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

	branch := extractBranch(id)

	return advancement.LocalizedAdvancement{ID: id, Title: e.Title, Description: e.Description, Icon: e.Icon, Difficulty: e.Difficulty, Branch: branch}, true
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

var branchNames = map[string]string{
	"story":      "Minecraft",
	"nether":     "Незер",
	"end":        "Энд",
	"adventure":  "Приключения",
	"husbandry":  "Сельское хозяйство",
}

func extractBranch(id string) string {
	after, ok := strings.CutPrefix(id, "minecraft:")
	if !ok {
		return ""
	}
	before, _, _ := strings.Cut(after, "/")
	name, ok := branchNames[before]
	if !ok {
		return before
	}
	return name
}

type entry struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Difficulty  string `json:"difficulty"`
}
