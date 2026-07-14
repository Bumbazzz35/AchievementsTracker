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
