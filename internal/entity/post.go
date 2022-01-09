package entity

import "time"

// Post is a short data about post in blog, without content.
// Content is available at the URL.
type Post struct {
	Title   string    `json:"title" bson:"title"`
	Date    time.Time `json:"date" bson:"date"`
	Author  string    `json:"author" bson:"author"`
	Summary string    `json:"summary" bson:"summary"`
	URL     string    `json:"url" bson:"url"`
}
