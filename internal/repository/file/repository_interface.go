package file

import (
	"context"

	"github.com/Bumbazzz35/AchievementTracker/internal/domain/advancement"
)

type PlayerAdvancementRepository interface {
	FindLatestWorldAdvancementFile(ctx context.Context) (string, error)
	ReadAdvancements(ctx context.Context, filePath string) (map[string]advancement.Advancement, error)
}
