package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// ListTurmasUseCase lista as turmas cadastradas
type ListTurmasUseCase struct{ repo repositories.TurmaRepository }

func NewListTurmasUseCase(repo repositories.TurmaRepository) *ListTurmasUseCase {
	return &ListTurmasUseCase{repo: repo}
}

// busca todas as turmas
func (u *ListTurmasUseCase) Execute(ctx context.Context) ([]*entities.Turma, error) {
	return u.repo.List(ctx)
}

// CreateTurmaUseCase cria uma turma nova
type CreateTurmaUseCase struct{ repo repositories.TurmaRepository }

func NewCreateTurmaUseCase(repo repositories.TurmaRepository) *CreateTurmaUseCase {
	return &CreateTurmaUseCase{repo: repo}
}

// exige nome antes de criar a turma
func (u *CreateTurmaUseCase) Execute(ctx context.Context, t *entities.Turma) error {
	if t.Name == "" {
		return errors.New("nome da turma é obrigatório")
	}
	return u.repo.Create(ctx, t)
}

// UpdateTurmaUseCase atualiza uma turma existente
type UpdateTurmaUseCase struct{ repo repositories.TurmaRepository }

func NewUpdateTurmaUseCase(repo repositories.TurmaRepository) *UpdateTurmaUseCase {
	return &UpdateTurmaUseCase{repo: repo}
}

// exige nome antes de salvar a atualizacao
func (u *UpdateTurmaUseCase) Execute(ctx context.Context, t *entities.Turma) error {
	if t.Name == "" {
		return errors.New("nome da turma é obrigatório")
	}
	return u.repo.Update(ctx, t)
}

// DeleteTurmaUseCase remove uma turma
type DeleteTurmaUseCase struct{ repo repositories.TurmaRepository }

func NewDeleteTurmaUseCase(repo repositories.TurmaRepository) *DeleteTurmaUseCase {
	return &DeleteTurmaUseCase{repo: repo}
}

// apaga a turma pelo id
func (u *DeleteTurmaUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
