package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

type BookmarkPostUseCase struct {
	repo ports.PostPort
}

func NewBookmarkPostUseCase(repo ports.PostPort) *BookmarkPostUseCase {
	return &BookmarkPostUseCase{repo: repo}
}

func (u *BookmarkPostUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.BookmarkPost(ctx, id)
}
