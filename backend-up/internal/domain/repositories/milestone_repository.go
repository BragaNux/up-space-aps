package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// MilestoneRepository cuida das conquistas/marcos de desenvolvimento dos alunos
type MilestoneRepository interface {
	List(ctx context.Context, studentID int64) ([]*entities.Milestone, error) // lista os marcos de um aluno
	GetByID(ctx context.Context, id int64) (*entities.Milestone, error)      // busca um marco especifico
	Create(ctx context.Context, milestone *entities.Milestone) error        // cria um marco novo
	Update(ctx context.Context, milestone *entities.Milestone) error        // atualiza um marco existente
	Delete(ctx context.Context, id int64) error                             // remove um marco
}
