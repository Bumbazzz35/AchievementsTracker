package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/internal/domain/advancement"
	"github.com/Bumbazzz35/AchievementTracker/internal/util/logger"
)

type RendererType string

const (
	RendererPlain RendererType = "plain"
	RendererColor RendererType = "color"
	RendererJSON  RendererType = "json"
)

type renderer interface {
	RenderStartup(event advancement.AdvancementEvent)
	RenderNewAdvancement(event advancement.AdvancementEvent)
	RenderProgressUpdate(event advancement.AdvancementEvent)
}

type Handler struct {
	renderer renderer
	logger   logger.Logger
}

func NewHandler(rt RendererType, logger logger.Logger, output io.Writer) *Handler {
	var r renderer

	switch rt {
	case RendererPlain:
		r = &plainRenderer{w: output}
	case RendererColor:
		r = &colorRenderer{w: output}
	case RendererJSON:
		r = &jsonRenderer{w: output}
	default:
		logger.Warnf("unknown renderer type %q, fallback to plain", rt)
		r = &plainRenderer{w: output}
	}

	return &Handler{renderer: r, logger: logger}
}

func (h *Handler) Start(ctx context.Context, eventCh <-chan advancement.AdvancementEvent) {
	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				h.logger.Info("event channel closed, stopping handler")
				return
			}

			switch event.Type {
			case advancement.EventTypeStartup:
				h.renderer.RenderStartup(event)
			case advancement.EventTypeNewAdvancement:
				h.renderer.RenderNewAdvancement(event)
			case advancement.EventTypeProgressUpdate:
				h.renderer.RenderProgressUpdate(event)
			}
		case <-ctx.Done():
			h.logger.Info("handler stopped by context")
			return
		}
	}
}

type plainRenderer struct {
	w io.Writer
}

func (r *plainRenderer) RenderStartup(event advancement.AdvancementEvent) {
	fmt.Fprintf(r.w, "\n=== Отслеживание мира %q начато! [%d/%d] ===\n", event.WorldName, event.Progress.Current, event.Progress.Total)
}

func (r *plainRenderer) RenderNewAdvancement(event advancement.AdvancementEvent) {
	fmt.Fprintf(r.w, "\n=== [%s] Новое достижение! [%d/%d] ===\n", event.WorldName, event.Progress.Current, event.Progress.Total)
	fmt.Fprintf(r.w, "  Название:      %s\n", event.Advancement.Title)
	fmt.Fprintf(r.w, "  Описание:      %s\n", event.Advancement.Description)
	fmt.Fprintf(r.w, "  Сложность:     %s\n", event.Advancement.Difficulty)
	fmt.Fprintf(r.w, "  Системный ID:  %s\n", event.Advancement.ID)
	fmt.Fprintf(r.w, "  Время:         %s\n", event.Timestamp.Format("02.01.2006 15:04:05"))
	fmt.Fprintln(r.w, strings.Repeat("=", 40))
}

func (r *plainRenderer) RenderProgressUpdate(event advancement.AdvancementEvent) {
	fmt.Fprintf(r.w, "\n--- Прогресс: [%d/%d] ---\n", event.Progress.Current, event.Progress.Total)
	if event.Advancement.Title != "" {
		fmt.Fprintf(r.w, "  %s\n", event.Advancement.Title)
	}
}

const (
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiBold   = "\033[1m"
	ansiReset  = "\033[0m"
)

type colorRenderer struct {
	w io.Writer
}

func (r *colorRenderer) RenderStartup(event advancement.AdvancementEvent) {
	fmt.Fprintf(r.w, "\n%s%s════════════════════════════════════════%s\n", ansiBold, ansiGreen, ansiReset)
	fmt.Fprintf(r.w, "%s%s 🚀 Отслеживание мира %q начато! [%d/%d]%s\n", ansiBold, ansiGreen, event.WorldName, event.Progress.Current, event.Progress.Total, ansiReset)
	fmt.Fprintf(r.w, "%s%s════════════════════════════════════════%s\n", ansiBold, ansiGreen, ansiReset)
}

func (r *colorRenderer) RenderNewAdvancement(event advancement.AdvancementEvent) {
	fmt.Fprintf(r.w, "\n%s%s════════════════════════════════════════%s\n", ansiBold, ansiGreen, ansiReset)
	fmt.Fprintf(r.w, "%s%s 🏆 [%s] Новое достижение! [%d/%d]%s\n", ansiBold, ansiGreen, event.WorldName, event.Progress.Current, event.Progress.Total, ansiReset)
	fmt.Fprintf(r.w, "%sНазвание:%s      %s%s%s\n", ansiYellow, ansiReset, ansiCyan, event.Advancement.Title, ansiReset)
	fmt.Fprintf(r.w, "%sОписание:%s      %s\n", ansiYellow, ansiReset, event.Advancement.Description)
	fmt.Fprintf(r.w, "%sСложность:%s     %s\n", ansiYellow, ansiReset, event.Advancement.Difficulty)
	fmt.Fprintf(r.w, "%sСистемный ID:%s  %s\n", ansiYellow, ansiReset, event.Advancement.ID)
	fmt.Fprintf(r.w, "%sВремя:%s         %s\n", ansiYellow, ansiReset, event.Timestamp.Format("02.01.2006 15:04:05"))
	fmt.Fprintf(r.w, "%s%s════════════════════════════════════════%s\n", ansiBold, ansiGreen, ansiReset)
}

func (r *colorRenderer) RenderProgressUpdate(event advancement.AdvancementEvent) {
	fmt.Fprintf(r.w, "\n%s--- Прогресс: [%d/%d] ---%s\n", ansiYellow, event.Progress.Current, event.Progress.Total, ansiReset)
	if event.Advancement.Title != "" {
		fmt.Fprintf(r.w, "%s%s%s%s\n", ansiCyan, ansiBold, event.Advancement.Title, ansiReset)
	}
}

type jsonRenderer struct {
	w io.Writer
}

func (r *jsonRenderer) RenderStartup(event advancement.AdvancementEvent) {
	je := jsonEvent{
		Type:      "startup",
		Progress:  event.Progress,
		WorldName: event.WorldName,
		Timestamp: event.Timestamp.Format(time.RFC3339),
	}

	data, _ := json.Marshal(je)
	fmt.Fprintln(r.w, string(data))
}

type jsonEvent struct {
	Type        string                          `json:"type"`
	Advancement jsonAdvancement                 `json:"advancement"`
	Progress    advancement.AdvancementProgress `json:"progress"`
	WorldName   string                          `json:"world_name"`
	Timestamp   string                          `json:"timestamp"`
}

type jsonAdvancement struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Difficulty  string `json:"difficulty,omitempty"`
}

func toJSONAdvancement(a advancement.LocalizedAdvancement) jsonAdvancement {
	return jsonAdvancement{
		ID:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		Icon:        a.Icon,
			Difficulty:  a.Difficulty,
	}
}

func (r *jsonRenderer) RenderNewAdvancement(event advancement.AdvancementEvent) {
	je := jsonEvent{
		Type:        "new_advancement",
		Advancement: toJSONAdvancement(event.Advancement),
		Progress:    event.Progress,
		WorldName:   event.WorldName,
		Timestamp:   event.Timestamp.Format(time.RFC3339),
	}

	data, _ := json.Marshal(je)
	fmt.Fprintln(r.w, string(data))
}

func (r *jsonRenderer) RenderProgressUpdate(event advancement.AdvancementEvent) {
	je := jsonEvent{
		Type:        "progress_update",
		Advancement: toJSONAdvancement(event.Advancement),
		Progress:    event.Progress,
		WorldName:   event.WorldName,
		Timestamp:   event.Timestamp.Format(time.RFC3339),
	}

	data, _ := json.Marshal(je)
	fmt.Fprintln(r.w, string(data))
}
