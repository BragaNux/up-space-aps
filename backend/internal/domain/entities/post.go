package entities

import "time"

type Post struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	PedagogicalNote string    `json:"pedagogical_note"`
	Likes           int64     `json:"likes"`
	Bookmarks       int64     `json:"bookmarks"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
