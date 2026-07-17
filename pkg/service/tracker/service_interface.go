package tracker

import (
	"context"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
)

type AchievementTrackerService interface {
	StartTracking(ctx context.Context) (<-chan advancement.AdvancementEvent, error)
	StopTracking()
	GetFullWorldState(ctx context.Context) (*advancement.FullWorldState, error)
}
