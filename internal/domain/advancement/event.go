package advancement

import "time"

type AdvancementEvent struct {
	Type        AdvancementEventType
	Advancement LocalizedAdvancement
	Progress    AdvancementProgress
	Timestamp   time.Time
	WorldName   string
}

type AdvancementEventType string

const (
	EventTypeStartup        AdvancementEventType = "startup"
	EventTypeNewAdvancement AdvancementEventType = "new_advancement"
	EventTypeProgressUpdate AdvancementEventType = "progress_update"
)

type AdvancementProgress struct {
	Current int
	Total   int
}
