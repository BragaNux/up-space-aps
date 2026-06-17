package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// LikePostUseCase da uma curtida num post
type LikePostUseCase struct {
	repo ports.PostPort
}

func NewLikePostUseCase(repo ports.PostPort) *LikePostUseCase {
	return &LikePostUseCase{repo: repo}
}

// soma a curtida e devolve o total atualizado
func (u *LikePostUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.LikePost(ctx, id)
}
