package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestReadAdvancements(t *testing.T) {
	tmpDir := t.TempDir()
	jsonContent := `{
		"minecraft:story/mine_stone": {
			"criteria": { "get_stone": "2026-06-26 17:12:58 +0300" },
			"done": true
		},
		"minecraft:story/root": {
			"done": false
		},
		"minecraft:recipes/misc/stick": {
			"criteria": {
			  "has_planks": "2026-07-13 15:58:30 +0300"
			},
			"done": true
		},
		"DataVersion": 4903
	}`

	filePath := filepath.Join(tmpDir, "test.json")
	if err := os.WriteFile(filePath, []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	repo := &fileRepository{}
	advancements, err := repo.ReadAdvancements(context.Background(), filePath)
	if err != nil {
		t.Fatalf("ReadAdvancements failed: %v", err)
	}

	if len(advancements) != 3 {
		t.Fatalf("expected 3 advancements, got %d", len(advancements))
	}

	stone, ok := advancements["minecraft:story/mine_stone"]
	if !ok {
		t.Fatal("expected minecraft:story/mine_stone to be present")
	}
	if stone.ID != "minecraft:story/mine_stone" {
		t.Fatalf("expected ID to be minecraft:story/mine_stone, got %q", stone.ID)
	}
	if !stone.Done {
		t.Fatal("expected stone to be done")
	}
	if len(stone.Criteria) != 1 {
		t.Fatalf("expected 1 criterion, got %d", len(stone.Criteria))
	}

	root, ok := advancements["minecraft:story/root"]
	if !ok {
		t.Fatal("expected minecraft:story/root to be present")
	}
	if root.ID != "minecraft:story/root" {
		t.Fatalf("expected ID to be minecraft:story/root, got %q", root.ID)
	}
	if root.Done {
		t.Fatal("expected root to be not done")
	}
	if len(root.Criteria) != 0 {
		t.Fatalf("expected 0 criteria, got %d", len(root.Criteria))
	}

	stickRecipe, ok := advancements["minecraft:recipes/misc/stick"]
	if !ok {
		t.Fatal("expected minecraft:recipes/misc/stick to be present")
	}
	if stickRecipe.ID != "minecraft:recipes/misc/stick" {
		t.Fatalf("expected ID to be minecraft:recipes/misc/stick, got %q", root.ID)
	}
	if !stickRecipe.Done {
		t.Fatal("expected stone to be done")
	}
	if len(stickRecipe.Criteria) != 1 {
		t.Fatalf("expected 1 criteria, got %d", len(root.Criteria))
	}
}
