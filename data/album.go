package data

type Album struct {
	User   string   `json:"user"`
	Albums []string `json:"albums"`
}
