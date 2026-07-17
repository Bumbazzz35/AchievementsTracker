package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
	"github.com/Bumbazzz35/AchievementTracker/pkg/util/logger"
)

type fileRepository struct {
	minecraftPath string
	logger        logger.Logger
}

func (repo *fileRepository) FindLatestWorldAdvancementFile(ctx context.Context) (string, error) {
	savesDir := filepath.Join(repo.minecraftPath, "saves")

	latestWorld, err := findLatestModified(savesDir, repo.logger)
	if err != nil {
		return "", fmt.Errorf("can't read folder with worlds - %q: %w", savesDir, err)
	}

	advancementsDir := filepath.Join(savesDir, latestWorld, "players", "advancements")

	latestPlayer, err := findLatestModified(advancementsDir, repo.logger)
	if err != nil {
		return "", fmt.Errorf("can't read folder with advancements - %q: %w", advancementsDir, err)
	}

	return filepath.Join(advancementsDir, latestPlayer), nil
}

func (repo *fileRepository) ReadAdvancements(ctx context.Context, filePath string) (map[string]advancement.Advancement, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't read file with advancements - %q: %w", filePath, err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed when parsing %q: %w", filePath, err)
	}

	advancements := make(map[string]advancement.Advancement, len(raw))
	for id, msg := range raw {
		var adv advancement.Advancement
		if err := json.Unmarshal(msg, &adv); err != nil {
			// Cannot parse "DataVersion": 4903
			// repo.logger.Warnf("cannot parse %q - %q with error: %v", id, msg, err)
			continue
		}
		adv.ID = id
		advancements[id] = adv
	}

	return advancements, nil
}

func findLatestModified(dir string, logger logger.Logger) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var latestName string
	var maxTime time.Time

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			logger.Warnf("Can't read info about %q: %v", file.Name(), err)
			continue
		}

		if latestName == "" || info.ModTime().After(maxTime) {
			maxTime = info.ModTime()
			latestName = file.Name()
		}
	}

	if latestName == "" {
		return "", fmt.Errorf("no files found in %s", dir)
	}

	return latestName, nil
}

func NewFileRepository(minecraftPath string, logger logger.Logger) *fileRepository {
	return &fileRepository{
		minecraftPath: minecraftPath,
		logger:        logger,
	}
}
