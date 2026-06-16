package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

type GetStudentUseCase struct {
	repo ports.StudentPort
}

func NewGetStudentUseCase(repo ports.StudentPort) *GetStudentUseCase {
	return &GetStudentUseCase{repo: repo}
}

func (u *GetStudentUseCase) Execute(ctx context.Context) (*entities.Student, error) {
	return u.repo.GetActive(ctx)
}
