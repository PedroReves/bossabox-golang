package model

type Tool struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
