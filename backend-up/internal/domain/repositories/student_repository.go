package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// StudentRepository cuida do cadastro e presenca dos alunos
type StudentRepository interface {
	GetActive(ctx context.Context) (*entities.Student, error)                                       // busca o aluno marcado como ativo no momento (uso de demo/single-student)
	GetByID(ctx context.Context, id int64) (*entities.Student, error)                                // busca um aluno especifico
	ListByGuardian(ctx context.Context, guardianUserID int64) ([]*entities.Student, error)           // lista os filhos vinculados a um responsavel
	List(ctx context.Context) ([]*entities.Student, error)                                           // lista todos os alunos
	Create(ctx context.Context, student *entities.Student) error                                     // cadastra um aluno novo
	Update(ctx context.Context, student *entities.Student) error                                     // atualiza dados de um aluno
	Delete(ctx context.Context, id int64) error                                                       // remove um aluno
	UpdatePresence(ctx context.Context, id int64, status string, checkInAt *string) error             // atualiza o status de presenca (e horario de entrada) do aluno
}
