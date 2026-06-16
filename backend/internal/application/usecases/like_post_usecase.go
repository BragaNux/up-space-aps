package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

type LikePostUseCase struct {
	repo ports.PostPort
}

func NewLikePostUseCase(repo ports.PostPort) *LikePostUseCase {
	return &LikePostUseCase{repo: repo}
}

func (u *LikePostUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.LikePost(ctx, id)
}
