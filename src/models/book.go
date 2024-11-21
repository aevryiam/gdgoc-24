package models

type Book struct {
	Id           string `json:"id" bson:"_id,omitempty"`
	Title        string `json:"title" bson:"title"`
	Author       string `json:"author" bson:"author"`
	Published_at string `json:"published_at" bson:"published_at"`
}
