package localization

import "testing"

func TestLocalizationRepository(t *testing.T) {
	repo, err := NewRepository()
	if err != nil {
		t.Fatalf("NewRepository failed: %v", err)
	}

	_, ok := repo.GetLocalizedAdvancement("minecraft:story/mine_stone")
	if !ok {
		t.Errorf("expected existing advancement to be found: %q", "minecraft:story/mine_stone")
	}

	_, ok = repo.GetLocalizedAdvancement("minecraft:non_existing/id")
	if ok {
		t.Errorf("expected unknown ID to return false: %q", "minecraft:non_existing/id")
	}

	advancements := repo.GetAllAdvancementIDs()
	if len(advancements) == 0 {
		t.Error("expected non-empty ID list")
	}
}
