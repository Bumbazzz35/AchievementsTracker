package tracker

import (
	"context"
	"path/filepath"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/internal/domain/advancement"
	"github.com/Bumbazzz35/AchievementTracker/internal/repository/file"
	"github.com/Bumbazzz35/AchievementTracker/internal/repository/localization"
	"github.com/Bumbazzz35/AchievementTracker/internal/util/fs"
	"github.com/Bumbazzz35/AchievementTracker/internal/util/logger"
)

type trackerService struct {
	fileRepo         file.PlayerAdvancementRepository
	localizationRepo localization.LocalizationRepository
	watcher          *fs.Watcher
	eventCh          chan advancement.AdvancementEvent
	cancel           context.CancelFunc
	pollInterval     time.Duration
	logger           logger.Logger
}

func (ts *trackerService) StartTracking(ctx context.Context) (<-chan advancement.AdvancementEvent, error) {
	ctx, ts.cancel = context.WithCancel(ctx)

	filePath, err := ts.fileRepo.FindLatestWorldAdvancementFile(ctx)
	if err != nil {
		return nil, err
	}
	oldAdvancements, err := ts.fileRepo.ReadAdvancements(ctx, filePath)
	if err != nil {
		return nil, err
	}

	currentAdv := getCurrentAdvancements(oldAdvancements, ts.localizationRepo)
	totalAdv := len(ts.localizationRepo.GetAllAdvancementIDs())

	ts.eventCh <- advancement.AdvancementEvent{
		Type:      advancement.EventTypeStartup,
		Progress:  advancement.AdvancementProgress{Current: currentAdv, Total: totalAdv},
		Timestamp: time.Now(),
		WorldName: filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(filePath)))),
	}

	ts.watcher = fs.NewWatcher(filePath, ts.pollInterval)
	watchCh := ts.watcher.Start(ctx)

	go func() {
		for {
			select {
			case msg := <-watchCh:
				newAdvancements, err := ts.fileRepo.ReadAdvancements(ctx, msg)
				if err != nil {
					ts.logger.Errorf("read advancements failed: %v", err)
					return
				}

				for id, adv := range newAdvancements {
					if !adv.Done {
						continue
					}
					old, exists := oldAdvancements[id]
					if exists && old.Done {
						continue
					}
					localAdv, ok := ts.localizationRepo.GetLocalizedAdvancement(id)
					if !ok {
						continue
					}
					currentAdv := getCurrentAdvancements(newAdvancements, ts.localizationRepo)
					totalAdv := len(ts.localizationRepo.GetAllAdvancementIDs())

					event := advancement.AdvancementEvent{
						Type:        advancement.EventTypeNewAdvancement,
						Advancement: localAdv,
						Progress:    advancement.AdvancementProgress{Current: currentAdv, Total: totalAdv},
						Timestamp:   time.Now(),
						WorldName:   filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(filePath)))),
					}
					select {
					case ts.eventCh <- event:
					case <-ctx.Done():
						return
					}
				}
				oldAdvancements = newAdvancements
			case <-ctx.Done():
				return
			}
		}
	}()

	return ts.eventCh, nil
}

func (ts *trackerService) StopTracking() {
	if ts.cancel != nil {
		ts.cancel()
	}
}

func NewTrackerService(fileRepo file.PlayerAdvancementRepository,
	localizationRepo localization.LocalizationRepository,
	pollInterval time.Duration,
	logger logger.Logger) *trackerService {
	return &trackerService{
		fileRepo:         fileRepo,
		localizationRepo: localizationRepo,
		watcher:          nil,
		pollInterval:     pollInterval,
		eventCh:          make(chan advancement.AdvancementEvent, 1),
		logger:           logger,
	}
}

func getCurrentAdvancements(newAdvancements map[string]advancement.Advancement, localizationRepo localization.LocalizationRepository) int {
	currentAdv := 0
	for id, adv := range newAdvancements {
		if !adv.Done {
			continue
		}

		// skip non-advancements e.g. recipes, etc
		if _, ok := localizationRepo.GetLocalizedAdvancement(id); !ok {
			continue
		}

		currentAdv++
	}

	return currentAdv
}
