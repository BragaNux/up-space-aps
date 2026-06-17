package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// UnlikePostUseCase tira a curtida de um post
type UnlikePostUseCase struct {
	repo ports.PostPort
}

func NewUnlikePostUseCase(repo ports.PostPort) *UnlikePostUseCase {
	return &UnlikePostUseCase{repo: repo}
}

// tira a curtida e devolve o total atualizado
func (u *UnlikePostUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.UnlikePost(ctx, id)
}
