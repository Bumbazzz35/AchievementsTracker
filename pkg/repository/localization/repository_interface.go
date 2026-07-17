package localization

import "github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"

type LocalizationRepository interface {
	GetLocalizedAdvancement(id string) (advancement.LocalizedAdvancement, bool)
	GetAllAdvancementIDs() []string
}
