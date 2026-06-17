package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// ListUsersByRoleUseCase lista usuarios filtrando por papel (profissional ou responsavel)
type ListUsersByRoleUseCase struct {
	repo repositories.UserRepository
}

func NewListUsersByRoleUseCase(repo repositories.UserRepository) *ListUsersByRoleUseCase {
	return &ListUsersByRoleUseCase{repo: repo}
}

// confere se o papel pedido e valido antes de listar
func (u *ListUsersByRoleUseCase) Execute(ctx context.Context, role string) ([]*entities.User, error) {
	if role != "profissional" && role != "responsavel" {
		return nil, errors.New("papel inválido")
	}
	return u.repo.ListByRole(ctx, role)
}
