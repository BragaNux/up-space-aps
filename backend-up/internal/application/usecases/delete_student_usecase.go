package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

// DeleteStudentUseCase remove o cadastro de um aluno
type DeleteStudentUseCase struct {
	repo ports.StudentPort
}

func NewDeleteStudentUseCase(repo ports.StudentPort) *DeleteStudentUseCase {
	return &DeleteStudentUseCase{repo: repo}
}

// apaga o aluno pelo id
func (u *DeleteStudentUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
