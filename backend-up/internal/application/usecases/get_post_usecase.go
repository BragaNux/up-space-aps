package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// GetPostUseCase busca um post especifico do feed
type GetPostUseCase struct {
	repo ports.PostPort
}

func NewGetPostUseCase(repo ports.PostPort) *GetPostUseCase {
	return &GetPostUseCase{repo: repo}
}

// busca o post pelo id
func (u *GetPostUseCase) Execute(ctx context.Context, id int64) (*entities.Post, error) {
	return u.repo.GetByID(ctx, id)
}
