package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// ListPostsUseCase lista os posts do feed (de um aluno especifico, se informado)
type ListPostsUseCase struct {
	repo ports.PostPort
}

func NewListPostsUseCase(repo ports.PostPort) *ListPostsUseCase {
	return &ListPostsUseCase{repo: repo}
}

// busca os posts, filtrando por aluno quando o id for informado
func (u *ListPostsUseCase) Execute(ctx context.Context, studentID int64) ([]*entities.Post, error) {
	return u.repo.ListPosts(ctx, studentID)
}
