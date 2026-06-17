package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// GuardianRepository cuida dos responsaveis autorizados a buscar cada aluno
type GuardianRepository interface {
	ListByStudent(ctx context.Context, studentID int64) ([]*entities.Guardian, error) // lista os responsaveis de um aluno
	GetByID(ctx context.Context, id int64) (*entities.Guardian, error)               // busca um responsavel especifico
	Create(ctx context.Context, guardian *entities.Guardian) error                   // cadastra um responsavel novo
	Update(ctx context.Context, guardian *entities.Guardian) error                   // atualiza dados de um responsavel
	Delete(ctx context.Context, id int64) error                                      // remove um responsavel
}
