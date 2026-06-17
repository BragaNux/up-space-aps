package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// categorias de marco que o sistema aceita
var validMilestoneCategories = map[string]bool{
	"Motor": true, "Linguagem": true, "Social": true, "Cognitivo": true,
}

// ListMilestonesUseCase lista os marcos de desenvolvimento de um aluno
type ListMilestonesUseCase struct {
	repo repositories.MilestoneRepository
}

func NewListMilestonesUseCase(repo repositories.MilestoneRepository) *ListMilestonesUseCase {
	return &ListMilestonesUseCase{repo: repo}
}

// busca os marcos daquele aluno
func (u *ListMilestonesUseCase) Execute(ctx context.Context, studentID int64) ([]*entities.Milestone, error) {
	return u.repo.List(ctx, studentID)
}

// CreateMilestoneUseCase registra um marco novo pra um aluno
type CreateMilestoneUseCase struct {
	repo repositories.MilestoneRepository
}

func NewCreateMilestoneUseCase(repo repositories.MilestoneRepository) *CreateMilestoneUseCase {
	return &CreateMilestoneUseCase{repo: repo}
}

// confere aluno, titulo e categoria antes de criar o marco
func (u *CreateMilestoneUseCase) Execute(ctx context.Context, m *entities.Milestone) error {
	if m.StudentID == 0 {
		return errors.New("id do aluno é obrigatório")
	}
	if m.Title == "" {
		return errors.New("título é obrigatório")
	}
	if !validMilestoneCategories[m.Category] {
		return errors.New("categoria inválida")
	}
	return u.repo.Create(ctx, m)
}

// UpdateMilestoneUseCase atualiza um marco existente
type UpdateMilestoneUseCase struct {
	repo repositories.MilestoneRepository
}

func NewUpdateMilestoneUseCase(repo repositories.MilestoneRepository) *UpdateMilestoneUseCase {
	return &UpdateMilestoneUseCase{repo: repo}
}

// confere titulo e categoria antes de salvar a atualizacao
func (u *UpdateMilestoneUseCase) Execute(ctx context.Context, m *entities.Milestone) error {
	if m.Title == "" {
		return errors.New("título é obrigatório")
	}
	if !validMilestoneCategories[m.Category] {
		return errors.New("categoria inválida")
	}
	return u.repo.Update(ctx, m)
}

// DeleteMilestoneUseCase remove um marco
type DeleteMilestoneUseCase struct {
	repo repositories.MilestoneRepository
}

func NewDeleteMilestoneUseCase(repo repositories.MilestoneRepository) *DeleteMilestoneUseCase {
	return &DeleteMilestoneUseCase{repo: repo}
}

// apaga o marco pelo id
func (u *DeleteMilestoneUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
