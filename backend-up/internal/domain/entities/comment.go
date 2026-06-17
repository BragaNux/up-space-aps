package entities

import "time"

// Comment e um comentario feito num post do feed
type Comment struct {
	ID         int64     `json:"id"`
	PostID     int64     `json:"post_id"`
	UserID     *int64    `json:"user_id,omitempty"`
	AuthorName string    `json:"author_name"`
	AvatarURL  string    `json:"avatar_url"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}
