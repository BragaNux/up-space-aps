package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// ListGuardiansUseCase lista os responsaveis de um aluno
type ListGuardiansUseCase struct {
	repo repositories.GuardianRepository
}

func NewListGuardiansUseCase(repo repositories.GuardianRepository) *ListGuardiansUseCase {
	return &ListGuardiansUseCase{repo: repo}
}

// busca os responsaveis daquele aluno
func (u *ListGuardiansUseCase) Execute(ctx context.Context, studentID int64) ([]*entities.Guardian, error) {
	return u.repo.ListByStudent(ctx, studentID)
}

// CreateGuardianUseCase cadastra um responsavel novo pra um aluno
type CreateGuardianUseCase struct {
	repo repositories.GuardianRepository
}

func NewCreateGuardianUseCase(repo repositories.GuardianRepository) *CreateGuardianUseCase {
	return &CreateGuardianUseCase{repo: repo}
}

// exige nome e parentesco antes de cadastrar
func (u *CreateGuardianUseCase) Execute(ctx context.Context, g *entities.Guardian) error {
	if g.Name == "" || g.Relation == "" {
		return errors.New("nome e parentesco são obrigatórios")
	}
	return u.repo.Create(ctx, g)
}

// UpdateGuardianUseCase atualiza os dados de um responsavel
type UpdateGuardianUseCase struct {
	repo repositories.GuardianRepository
}

func NewUpdateGuardianUseCase(repo repositories.GuardianRepository) *UpdateGuardianUseCase {
	return &UpdateGuardianUseCase{repo: repo}
}

// exige nome e parentesco antes de atualizar
func (u *UpdateGuardianUseCase) Execute(ctx context.Context, g *entities.Guardian) error {
	if g.Name == "" || g.Relation == "" {
		return errors.New("nome e parentesco são obrigatórios")
	}
	return u.repo.Update(ctx, g)
}

// DeleteGuardianUseCase remove um responsavel
type DeleteGuardianUseCase struct {
	repo repositories.GuardianRepository
}

func NewDeleteGuardianUseCase(repo repositories.GuardianRepository) *DeleteGuardianUseCase {
	return &DeleteGuardianUseCase{repo: repo}
}

// apaga o responsavel pelo id
func (u *DeleteGuardianUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
