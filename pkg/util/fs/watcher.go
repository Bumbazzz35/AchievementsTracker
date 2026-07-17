package fs

import (
	"context"
	"os"
	"time"
)

type Watcher struct {
	filePath    string
	interval    time.Duration
	lastModTime time.Time
	lastSize    int64 // Подстраховка от глюков lastModTime на винде
}

func NewWatcher(filePath string, interval time.Duration) *Watcher {
	return &Watcher{
		filePath:    filePath,
		interval:    interval,
		lastModTime: time.Time{},
		lastSize:    -1,
	}
}

func (w *Watcher) Start(ctx context.Context) <-chan string {
	ch := make(chan string, 1)

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-ticker.C:
				info, err := os.Stat(w.filePath)
				if err != nil {
					continue
				}

				if w.lastModTime == info.ModTime() && w.lastSize == info.Size() {
					continue
				}

				w.lastModTime = info.ModTime()
				w.lastSize = info.Size()
				ch <- w.filePath
			}
		}
	}()

	return ch
}
