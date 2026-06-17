package dto

// CreateCommentRequest e o corpo esperado no POST de criar comentario num post
type CreateCommentRequest struct {
	Text string `json:"text"`
}
