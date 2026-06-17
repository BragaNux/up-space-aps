package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// GetStudentUseCase busca o aluno "ativo" (usado no modo single-student/demo)
type GetStudentUseCase struct {
	repo ports.StudentPort
}

func NewGetStudentUseCase(repo ports.StudentPort) *GetStudentUseCase {
	return &GetStudentUseCase{repo: repo}
}

// busca o aluno marcado como ativo
func (u *GetStudentUseCase) Execute(ctx context.Context) (*entities.Student, error) {
	return u.repo.GetActive(ctx)
}
