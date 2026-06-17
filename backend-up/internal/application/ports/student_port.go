package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// StudentPort e a porta que os usecases de aluno usam pra falar com o repositorio
type StudentPort interface {
	GetActive(ctx context.Context) (*entities.Student, error)
	GetByID(ctx context.Context, id int64) (*entities.Student, error)
	ListByGuardian(ctx context.Context, guardianUserID int64) ([]*entities.Student, error)
	List(ctx context.Context) ([]*entities.Student, error)
	Create(ctx context.Context, student *entities.Student) error
	Update(ctx context.Context, student *entities.Student) error
	Delete(ctx context.Context, id int64) error
	UpdatePresence(ctx context.Context, id int64, status string, checkInAt *string) error
}
