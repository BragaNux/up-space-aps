package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// ListStudentsUseCase lista todos os alunos cadastrados
type ListStudentsUseCase struct {
	repo ports.StudentPort
}

func NewListStudentsUseCase(repo ports.StudentPort) *ListStudentsUseCase {
	return &ListStudentsUseCase{repo: repo}
}

// busca todos os alunos
func (u *ListStudentsUseCase) Execute(ctx context.Context) ([]*entities.Student, error) {
	return u.repo.List(ctx)
}
