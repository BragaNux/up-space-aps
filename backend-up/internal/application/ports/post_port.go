package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// PostPort e a porta que os usecases de post usam pra falar com o repositorio
type PostPort interface {
	ListPosts(ctx context.Context, studentID int64) ([]*entities.Post, error)
	GetByID(ctx context.Context, id int64) (*entities.Post, error)
	CreatePost(ctx context.Context, post *entities.Post) error
	UpdatePost(ctx context.Context, post *entities.Post) error
	DeletePost(ctx context.Context, id int64) error
	LikePost(ctx context.Context, id int64) (int64, error)
	UnlikePost(ctx context.Context, id int64) (int64, error)
	BookmarkPost(ctx context.Context, id int64) (int64, error)
	UnbookmarkPost(ctx context.Context, id int64) (int64, error)
}
