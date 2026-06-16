package dto

type CreatePostRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	PedagogicalNote string `json:"pedagogical_note"`
}
