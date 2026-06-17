package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// CommentRepository cuida dos comentarios feitos nos posts do feed
type CommentRepository interface {
	ListByPost(ctx context.Context, postID int64) ([]*entities.Comment, error) // lista os comentarios de um post
	GetByID(ctx context.Context, id int64) (*entities.Comment, error)         // busca um comentario especifico
	Create(ctx context.Context, comment *entities.Comment) error             // cria um comentario novo
	Delete(ctx context.Context, id int64) error                              // remove um comentario
}
