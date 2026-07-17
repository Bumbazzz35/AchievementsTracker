package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
	"github.com/Bumbazzz35/AchievementTracker/pkg/repository/file"
	"github.com/Bumbazzz35/AchievementTracker/pkg/repository/localization"
	servicetracker "github.com/Bumbazzz35/AchievementTracker/pkg/service/tracker"
	"github.com/Bumbazzz35/AchievementTracker/pkg/util/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppConfig struct {
	MinecraftPath string `json:"minecraftPath"`
	PollInterval  int    `json:"pollInterval"`
}

type FrontendEvent struct {
	Type            string   `json:"type"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Icon            string   `json:"icon"`
	Difficulty      string   `json:"difficulty"`
	Branch          string   `json:"branch"`
	ProgressCurrent int      `json:"progressCurrent"`
	ProgressTotal   int      `json:"progressTotal"`
	WorldName       string   `json:"worldName"`
	CriteriaUpdates []string `json:"criteriaUpdates"`
}

type App struct {
	ctx    context.Context
	svc    servicetracker.AchievementTrackerService
	log    logger.Logger
	config AppConfig
	stop   context.CancelFunc
}

func NewApp() *App {
	log := logger.New("info", os.Stderr)

	config := loadConfig()

	minecraftPath := config.MinecraftPath
	if minecraftPath == "" {
		minecraftPath = filepath.Join(os.Getenv("APPDATA"), ".minecraft")
	}

	fileRepo := file.NewFileRepository(minecraftPath, log)

	locRepo, err := localization.NewRepository()
	if err != nil {
		log.Errorf("localization init failed: %v", err)
		return &App{log: log, config: config}
	}

	pollInterval := time.Duration(config.PollInterval) * time.Second
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	svc := servicetracker.NewTrackerService(fileRepo, locRepo, pollInterval, log)

	return &App{svc: svc, log: log, config: config}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	a.StopWatching()
}

func (a *App) GetFullWorldState() (*advancement.FullWorldState, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.GetFullWorldState(a.ctx)
}

func (a *App) StartWatching() error {
	if a.svc == nil {
		return fmt.Errorf("service not initialized")
	}

	eventCh, err := a.svc.StartTracking(a.ctx)
	if err != nil {
		a.log.Errorf("StartTracking failed: %v", err)
		return err
	}

	a.log.Info("StartTracking started, waiting for events...")

	ctx, cancel := context.WithCancel(a.ctx)
	a.stop = cancel

	go func() {
		defer cancel()
		for {
			select {
			case event, ok := <-eventCh:
				if !ok {
					a.log.Info("event channel closed")
					return
				}

				fe := FrontendEvent{
					Type:            string(event.Type),
					Title:           event.Advancement.Title,
					Description:     event.Advancement.Description,
					Icon:            event.Advancement.Icon,
					Difficulty:      event.Advancement.Difficulty,
					Branch:          event.Advancement.Branch,
					ProgressCurrent: event.Progress.Current,
					ProgressTotal:   event.Progress.Total,
					WorldName:       event.WorldName,
					CriteriaUpdates: event.CriteriaUpdates,
				}

				a.log.Infof("emitting event: type=%s title=%q", fe.Type, fe.Title)
				runtime.EventsEmit(a.ctx, "advancement-event", fe)
			case <-ctx.Done():
				a.log.Info("event processing stopped")
				return
			}
		}
	}()

	return nil
}

func (a *App) StopWatching() {
	if a.stop != nil {
		a.stop()
	}
	if a.svc != nil {
		a.svc.StopTracking()
	}
}

func (a *App) SaveConfig(cfg AppConfig) error {
	a.config = cfg
	return saveConfig(cfg)
}

func (a *App) LoadConfig() AppConfig {
	return a.config
}

func (a *App) SendTestEvent() {
	runtime.EventsEmit(a.ctx, "advancement-event", FrontendEvent{
		Type:            "startup",
		ProgressCurrent: 42,
		ProgressTotal:   126,
		WorldName:       "ТЕСТ",
	})
}

func configPath() string {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = "."
	}
	dir := filepath.Join(appData, "AchievementTracker")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "config.json")
}

func loadConfig() AppConfig {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return AppConfig{PollInterval: 5}
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{PollInterval: 5}
	}
	return cfg
}

func saveConfig(cfg AppConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}
