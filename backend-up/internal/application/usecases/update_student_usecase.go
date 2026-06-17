package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// UpdateStudentUseCase atualiza o cadastro de um aluno
type UpdateStudentUseCase struct {
	repo ports.StudentPort
}

func NewUpdateStudentUseCase(repo ports.StudentPort) *UpdateStudentUseCase {
	return &UpdateStudentUseCase{repo: repo}
}

// confere nome e foto antes de salvar a atualizacao
func (u *UpdateStudentUseCase) Execute(ctx context.Context, student *entities.Student) error {
	if student.Name == "" {
		return errors.New("nome é obrigatório")
	}
	if err := validateImage(student.PhotoURL); err != nil {
		return err
	}
	return u.repo.Update(ctx, student)
}
