package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type PostRepository interface {
	List(ctx context.Context) ([]*entities.Post, error)
	Create(ctx context.Context, post *entities.Post) error
	IncrementLikes(ctx context.Context, id int64) (int64, error)
	IncrementBookmarks(ctx context.Context, id int64) (int64, error)
}
