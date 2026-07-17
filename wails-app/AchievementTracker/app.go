package main

import (
	"context"
	"os"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
	"github.com/Bumbazzz35/AchievementTracker/pkg/repository/file"
	"github.com/Bumbazzz35/AchievementTracker/pkg/repository/localization"
	servicetracker "github.com/Bumbazzz35/AchievementTracker/pkg/service/tracker"
	"github.com/Bumbazzz35/AchievementTracker/pkg/util/logger"
)

type App struct {
	ctx context.Context
	svc servicetracker.AchievementTrackerService
	log logger.Logger
}

func NewApp() *App {
	log := logger.New("info", os.Stderr)

	minecraftPath := os.Getenv("APPDATA") + "\\.minecraft"

	fileRepo := file.NewFileRepository(minecraftPath, log)

	locRepo, err := localization.NewRepository()
	if err != nil {
		log.Errorf("localization init failed: %v", err)
		return &App{log: log}
	}

	svc := servicetracker.NewTrackerService(fileRepo, locRepo, 0, log)

	return &App{svc: svc, log: log}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetFullWorldState() (*advancement.FullWorldState, error) {
	if a.svc == nil {
		return nil, nil
	}
	return a.svc.GetFullWorldState(a.ctx)
}
