package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// BookmarkPostUseCase salva um post na lista de favoritos do usuario
type BookmarkPostUseCase struct {
	repo ports.PostPort
}

func NewBookmarkPostUseCase(repo ports.PostPort) *BookmarkPostUseCase {
	return &BookmarkPostUseCase{repo: repo}
}

// soma um "salvo" no post e devolve o total atualizado
func (u *BookmarkPostUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.BookmarkPost(ctx, id)
}
