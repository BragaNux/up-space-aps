package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// TurmaRepository cuida do cadastro das turmas/salas
type TurmaRepository interface {
	List(ctx context.Context) ([]*entities.Turma, error)        // lista todas as turmas
	GetByID(ctx context.Context, id int64) (*entities.Turma, error) // busca uma turma especifica
	Create(ctx context.Context, turma *entities.Turma) error    // cria uma turma nova
	Update(ctx context.Context, turma *entities.Turma) error    // atualiza uma turma existente
	Delete(ctx context.Context, id int64) error                 // remove uma turma
}
