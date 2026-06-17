package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// CreateStudentUseCase cadastra um aluno novo
type CreateStudentUseCase struct {
	repo ports.StudentPort
}

func NewCreateStudentUseCase(repo ports.StudentPort) *CreateStudentUseCase {
	return &CreateStudentUseCase{repo: repo}
}

// confere nome, data de nascimento e foto antes de cadastrar o aluno
func (u *CreateStudentUseCase) Execute(ctx context.Context, student *entities.Student) error {
	if student.Name == "" {
		return errors.New("nome é obrigatório")
	}
	if student.BirthDate == nil {
		return errors.New("data de nascimento é obrigatória")
	}
	if err := validateImage(student.PhotoURL); err != nil {
		return err
	}
	return u.repo.Create(ctx, student)
}
