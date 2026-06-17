package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// DeletePostUseCase remove um post do feed
type DeletePostUseCase struct {
	repo ports.PostPort
}

func NewDeletePostUseCase(repo ports.PostPort) *DeletePostUseCase {
	return &DeletePostUseCase{repo: repo}
}

// apaga o post pelo id
func (u *DeletePostUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.DeletePost(ctx, id)
}
