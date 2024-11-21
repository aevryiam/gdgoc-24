package models

type BookResponse struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Published_at string `json:"published_at"`
}
