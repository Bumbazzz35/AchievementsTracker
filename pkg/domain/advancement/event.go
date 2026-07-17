package advancement

import "time"

type AdvancementEvent struct {
	Type            AdvancementEventType
	Advancement     LocalizedAdvancement
	Progress        AdvancementProgress
	Timestamp       time.Time
	WorldName       string
	CriteriaUpdates []string // "minecraft:plains"
}

type AdvancementEventType string

const (
	EventTypeStartup        AdvancementEventType = "startup"
	EventTypeNewAdvancement AdvancementEventType = "new_advancement"
	EventTypeProgressUpdate AdvancementEventType = "progress_update"
	EventTypeCriteriaUpdate AdvancementEventType = "criteria_update"
)

type AdvancementProgress struct {
	Current int
	Total   int
}
