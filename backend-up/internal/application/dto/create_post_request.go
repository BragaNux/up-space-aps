package dto

// CreatePostRequest e o corpo esperado no POST /api/posts
type CreatePostRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	PedagogicalNote string `json:"pedagogical_note"`
	ImageURL        string `json:"image_url"`
	Visibility      string `json:"visibility"`
}
