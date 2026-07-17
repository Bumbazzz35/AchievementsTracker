package tracker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/pkg/delivery/cli"
	"github.com/Bumbazzz35/AchievementTracker/pkg/repository/file"
	"github.com/Bumbazzz35/AchievementTracker/pkg/repository/localization"
	servicetracker "github.com/Bumbazzz35/AchievementTracker/pkg/service/tracker"
	"github.com/Bumbazzz35/AchievementTracker/pkg/util/logger"
)

type App struct {
	logger        logger.Logger
	minecraftPath string
	pollInterval  time.Duration
	rendererType  cli.RendererType
}

func (a *App) Run() error {
	a.logger.Info("starting app")

	fileRepo := file.NewFileRepository(a.minecraftPath, a.logger)

	locRepo, err := localization.NewRepository()
	if err != nil {
		return fmt.Errorf("localization: %w", err)
	}

	svc := servicetracker.NewTrackerService(fileRepo, locRepo, a.pollInterval, a.logger)

	handler := cli.NewHandler(a.rendererType, a.logger, os.Stdout)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	eventCh, err := svc.StartTracking(ctx)
	if err != nil {
		return fmt.Errorf("start tracking: %w", err)
	}

	go handler.Start(ctx, eventCh)

	sig := <-sigCh
	a.logger.Infof("received signal %v, shutting down", sig)

	cancel()
	svc.StopTracking()
	return nil
}

func NewApp(cfg Config) *App {
	return &App{
		logger:        logger.New(cfg.LogLevel, cfg.LogOutput),
		minecraftPath: cfg.MinecraftPath,
		pollInterval:  cfg.PollInterval,
		rendererType:  cfg.RendererType,
	}
}

type Config struct {
	MinecraftPath string
	PollInterval  time.Duration
	RendererType  cli.RendererType
	LogLevel      string // "debug", "info", "warn", "error"
	LogOutput     io.Writer
}
