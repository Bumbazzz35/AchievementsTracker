package tracker

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/internal/domain/advancement"
	"github.com/Bumbazzz35/AchievementTracker/internal/repository/file"
	"github.com/Bumbazzz35/AchievementTracker/internal/repository/localization"
	"github.com/Bumbazzz35/AchievementTracker/internal/util/fs"
	"github.com/Bumbazzz35/AchievementTracker/internal/util/logger"
)

var branchNames = map[string]string{
	"story":     "Minecraft",
	"nether":    "Незер",
	"end":       "Энд",
	"adventure": "Приключения",
	"husbandry": "Сельское хозяйство",
}

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
					localAdv, ok := ts.localizationRepo.GetLocalizedAdvancement(id)
					// skip recipes and non-advancements
					if !ok {
						continue
					}

					oldAdv, exists := oldAdvancements[id]
					if exists && oldAdv.Done {
						continue
					}

					currentAdv := getCurrentAdvancements(newAdvancements, ts.localizationRepo)
					totalAdv := len(ts.localizationRepo.GetAllAdvancementIDs())
					worldName := filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(filePath))))

					if adv.Done {
						event := advancement.AdvancementEvent{
							Type:        advancement.EventTypeNewAdvancement,
							Advancement: localAdv,
							Progress:    advancement.AdvancementProgress{Current: currentAdv, Total: totalAdv},
							Timestamp:   time.Now(),
							WorldName:   worldName,
						}
						select {
						case ts.eventCh <- event:
						case <-ctx.Done():
							return
						}
					} else if exists {
						var updates []string
						for key := range adv.Criteria {
							if _, ok := oldAdv.Criteria[key]; !ok {
								updates = append(updates, key)
							}
						}

						if len(updates) > 0 {
							event := advancement.AdvancementEvent{
								Type:            advancement.EventTypeCriteriaUpdate,
								Advancement:     localAdv,
								Progress:        advancement.AdvancementProgress{Current: currentAdv, Total: totalAdv},
								Timestamp:       time.Now(),
								WorldName:       worldName,
								CriteriaUpdates: updates,
							}

							select {
							case ts.eventCh <- event:
							case <-ctx.Done():
								return
							}
						}
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

func (ts *trackerService) GetFullWorldState(ctx context.Context) (*advancement.FullWorldState, error) {
	filePath, err := ts.fileRepo.FindLatestWorldAdvancementFile(ctx)
	if err != nil {
		return nil, err
	}

	branchMap := make(map[string]*advancement.BranchSnapshot, 5)

	advancements, err := ts.fileRepo.ReadAdvancements(ctx, filePath)
	if err != nil {
		return nil, err
	}

	for _, id := range ts.localizationRepo.GetAllAdvancementIDs() {
		locAdv, _ := ts.localizationRepo.GetLocalizedAdvancement(id)

		adv, ok := advancements[locAdv.ID]

		var done bool
		var criteria map[string]bool

		if ok {
			done = adv.Done
			criteria = convertCriteria(adv.Criteria)
		}

		item := advancement.ItemWithCriteria{
			LocalizedAdvancement: locAdv,
			Done:                 done,
			Criteria:             criteria,
			IsBig:                advancement.BigAchievementIDs[locAdv.ID],
		}
		branchID := extractBranchID(item.ID)

		branchSnapshot, ok := branchMap[branchID]
		if ok {
			branchSnapshot.Items = append(branchSnapshot.Items, item)
		} else {
			branchMap[branchID] = &advancement.BranchSnapshot{
				ID:    branchID,
				Title: item.Branch,
				Items: []advancement.ItemWithCriteria{item},
			}
		}
	}

	branches := make([]advancement.BranchSnapshot, 0, len(branchMap))
	for _, v := range branchMap {
		doneCount := 0
		for _, item := range v.Items {
			if item.Done {
				doneCount++
			}
		}

		v.DoneCount = doneCount
		v.TotalCount = len(v.Items)
		branches = append(branches, *v)
	}

	return &advancement.FullWorldState{
		WorldName: filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(filePath)))),
		Branches:  branches,
	}, nil
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

func convertCriteria(criteria map[string]string) map[string]bool {
	m := make(map[string]bool, len(criteria))
	for key := range criteria {
		m[key] = true
	}

	return m
}

func extractBranchID(id string) string {
	after, ok := strings.CutPrefix(id, "minecraft:")
	if !ok {
		return ""
	}
	before, _, _ := strings.Cut(after, "/")
	return before
}
