package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

type ListPostsUseCase struct {
	repo ports.PostPort
}

func NewListPostsUseCase(repo ports.PostPort) *ListPostsUseCase {
	return &ListPostsUseCase{repo: repo}
}

func (u *ListPostsUseCase) Execute(ctx context.Context) ([]*entities.Post, error) {
	return u.repo.ListPosts(ctx)
}
