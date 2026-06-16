package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type PostPort interface {
	ListPosts(ctx context.Context) ([]*entities.Post, error)
	CreatePost(ctx context.Context, post *entities.Post) error
	LikePost(ctx context.Context, id int64) (int64, error)
	BookmarkPost(ctx context.Context, id int64) (int64, error)
}
