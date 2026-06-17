package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// UnbookmarkPostUseCase tira um post da lista de favoritos
type UnbookmarkPostUseCase struct {
	repo ports.PostPort
}

func NewUnbookmarkPostUseCase(repo ports.PostPort) *UnbookmarkPostUseCase {
	return &UnbookmarkPostUseCase{repo: repo}
}

// tira o "salvo" do post e devolve o total atualizado
func (u *UnbookmarkPostUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.UnbookmarkPost(ctx, id)
}
