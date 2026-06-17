package entities

import "time"

// Post e uma publicacao do feed pedagogico sobre um aluno (com foto, curtidas, comentarios etc)
type Post struct {
	ID              int64     `json:"id"`
	StudentID       int64     `json:"student_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	PedagogicalNote string    `json:"pedagogical_note"`
	ImageURL        string    `json:"image_url"`
	Likes           int64     `json:"likes"`
	Bookmarks       int64     `json:"bookmarks"`
	Visibility      string    `json:"visibility"`
	StudentName     string    `json:"student_name,omitempty"`
	TeacherName     string    `json:"teacher_name,omitempty"`
	TeacherAvatar   string    `json:"teacher_avatar,omitempty"`
	CommentCount    int64     `json:"comment_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
