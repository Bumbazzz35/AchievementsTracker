package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Bumbazzz35/AchievementTracker/pkg/domain/advancement"
)

func testEvent() advancement.AdvancementEvent {
	return advancement.AdvancementEvent{
		Type: advancement.EventTypeNewAdvancement,
		Advancement: advancement.LocalizedAdvancement{
			ID:          "minecraft:story/mine_stone",
			Title:       "Каменный век",
			Description: "Добудьте камень новой киркой",
			Icon:        "minecraft:wooden_pickaxe",
			Difficulty:  "Достижение",
		},
		Progress:  advancement.AdvancementProgress{Current: 3, Total: 126},
		WorldName: "Новый мир",
		Timestamp: time.Date(2026, 7, 13, 15, 30, 0, 0, time.UTC),
	}
}

func TestPlainRenderer_RenderStartup(t *testing.T) {
	var buf bytes.Buffer
	r := &plainRenderer{w: &buf}

	event := testEvent()
	event.Type = advancement.EventTypeStartup
	r.RenderStartup(event)

	out := buf.String()
	if !strings.Contains(out, "Отслеживание мира") {
		t.Error("expected 'Отслеживание мира'")
	}
	if !strings.Contains(out, "Новый мир") {
		t.Error("expected world name")
	}
	if !strings.Contains(out, "3/126") {
		t.Error("expected progress")
	}
}

func TestPlainRenderer_RenderNewAdvancement(t *testing.T) {
	var buf bytes.Buffer
	r := &plainRenderer{w: &buf}

	r.RenderNewAdvancement(testEvent())

	out := buf.String()
	if !strings.Contains(out, "Новое достижение") {
		t.Error("expected 'Новое достижение'")
	}
	if !strings.Contains(out, "Каменный век") {
		t.Error("expected title")
	}
	if !strings.Contains(out, "Добудьте камень новой киркой") {
		t.Error("expected description")
	}
	if !strings.Contains(out, "Новый мир") {
		t.Error("expected world name")
	}
}

func TestPlainRenderer_RenderProgressUpdate(t *testing.T) {
	var buf bytes.Buffer
	r := &plainRenderer{w: &buf}

	event := testEvent()
	event.Type = advancement.EventTypeProgressUpdate
	r.RenderProgressUpdate(event)

	out := buf.String()
	if !strings.Contains(out, "Прогресс") {
		t.Error("expected 'Прогресс'")
	}
	if !strings.Contains(out, "3/126") {
		t.Error("expected progress")
	}
}

func TestColorRenderer_RenderNewAdvancement(t *testing.T) {
	var buf bytes.Buffer
	r := &colorRenderer{w: &buf}

	r.RenderNewAdvancement(testEvent())

	out := buf.String()
	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes")
	}
	if !strings.Contains(out, "Новое достижение") {
		t.Error("expected 'Новое достижение'")
	}
}

func TestJSONRenderer_RenderStartup(t *testing.T) {
	var buf bytes.Buffer
	r := &jsonRenderer{w: &buf}

	event := testEvent()
	event.Type = advancement.EventTypeStartup
	r.RenderStartup(event)

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["type"] != "startup" {
		t.Errorf("expected type 'startup', got %q", parsed["type"])
	}
	if parsed["world_name"] != "Новый мир" {
		t.Errorf("expected world_name 'Новый мир', got %q", parsed["world_name"])
	}
}

func TestJSONRenderer_RenderNewAdvancement(t *testing.T) {
	var buf bytes.Buffer
	r := &jsonRenderer{w: &buf}

	r.RenderNewAdvancement(testEvent())

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["type"] != "new_advancement" {
		t.Errorf("expected type 'new_advancement', got %q", parsed["type"])
	}

	adv, ok := parsed["advancement"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'advancement' object")
	}
	if adv["title"] != "Каменный век" {
		t.Errorf("expected title 'Каменный век', got %q", adv["title"])
	}
	if adv["id"] != "minecraft:story/mine_stone" {
		t.Errorf("expected id 'minecraft:story/mine_stone', got %q", adv["id"])
	}
}

func TestJSONRenderer_RenderProgressUpdate(t *testing.T) {
	var buf bytes.Buffer
	r := &jsonRenderer{w: &buf}

	event := testEvent()
	event.Type = advancement.EventTypeProgressUpdate
	r.RenderProgressUpdate(event)

	var parsed map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if parsed["type"] != "progress_update" {
		t.Errorf("expected type 'progress_update', got %q", parsed["type"])
	}
}
