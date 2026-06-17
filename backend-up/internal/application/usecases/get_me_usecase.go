package usecases

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// GetMeUseCase busca os dados do usuario logado (rota /api/me)
type GetMeUseCase struct {
	repo repositories.UserRepository
}

func NewGetMeUseCase(repo repositories.UserRepository) *GetMeUseCase {
	return &GetMeUseCase{repo: repo}
}

// busca o usuario pelo id que veio do token
func (u *GetMeUseCase) Execute(ctx context.Context, userID int64) (*entities.User, error) {
	return u.repo.GetByID(ctx, userID)
}
