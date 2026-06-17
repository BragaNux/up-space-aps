package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// ListMyChildrenUseCase lista os filhos vinculados a um responsavel logado
type ListMyChildrenUseCase struct {
	repo ports.StudentPort
}

func NewListMyChildrenUseCase(repo ports.StudentPort) *ListMyChildrenUseCase {
	return &ListMyChildrenUseCase{repo: repo}
}

// busca os alunos vinculados aquele responsavel
func (u *ListMyChildrenUseCase) Execute(ctx context.Context, guardianUserID int64) ([]*entities.Student, error) {
	return u.repo.ListByGuardian(ctx, guardianUserID)
}
