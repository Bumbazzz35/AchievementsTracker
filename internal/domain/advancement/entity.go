package advancement

type Advancement struct {
	ID       string            `json:"-"`
	Done     bool              `json:"done"`
	Criteria map[string]string `json:"criteria"`
}

type LocalizedAdvancement struct {
	ID          string
	Title       string
	Description string
	Icon        string
	Difficulty  string
	Branch      string
}

type FullWorldState struct {
	WorldName string
	Branches  []BranchSnapshot
}

type BranchSnapshot struct {
	ID         string // "story", "nether", "end", "advancement", "husbandry"
	Title      string // "Minecraft", "Nether", "End", "Advancement", "Husbandry"
	Items      []ItemWithCriteria
	DoneCount  int
	TotalCount int
}

type ItemWithCriteria struct {
	LocalizedAdvancement
	Done     bool
	IsBig    bool
	Criteria map[string]bool // criteria_id -> true/false <--> "minecraft:plains": true/false
}
