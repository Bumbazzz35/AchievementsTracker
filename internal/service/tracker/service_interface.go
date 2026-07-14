package tracker

import (
	"context"

	"github.com/Bumbazzz35/AchievementTracker/internal/domain/advancement"
)

type AchievementTrackerService interface {
	StartTracking(ctx context.Context) (<-chan advancement.AdvancementEvent, error)
	StopTracking()
}
