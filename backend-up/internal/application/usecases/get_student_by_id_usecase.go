package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// GetStudentByIDUseCase busca um aluno especifico pelo id
type GetStudentByIDUseCase struct {
	repo ports.StudentPort
}

func NewGetStudentByIDUseCase(repo ports.StudentPort) *GetStudentByIDUseCase {
	return &GetStudentByIDUseCase{repo: repo}
}

// busca o aluno pelo id
func (u *GetStudentByIDUseCase) Execute(ctx context.Context, id int64) (*entities.Student, error) {
	return u.repo.GetByID(ctx, id)
}
