package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/pkg/app/tracker"
	"github.com/Bumbazzz35/AchievementTracker/pkg/delivery/cli"
)

func main() {
	var (
		minecraftPath = flag.String("minecraftPath", filepath.Join(os.Getenv("APPDATA"), ".minecraft"), "path to .minecraft folder")
		pollInterval  = flag.Duration("checkInterval", 5*time.Second, "interval to check advancement file")
		rendererType  = flag.String("renderer-type", "color", "output format: plain(default), color, json")
		logLevel      = flag.String("log-level", "info", "log level: debug, info, warn, error")
	)
	flag.Parse()

	cfg := tracker.Config{
		MinecraftPath: *minecraftPath,
		PollInterval:  *pollInterval,
		RendererType:  cli.RendererType(*rendererType),
		LogLevel:      *logLevel,
		LogOutput:     os.Stderr,
	}

	app := tracker.NewApp(cfg)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
