package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

type StudentPort interface {
	GetActive(ctx context.Context) (*entities.Student, error)
	UpdatePresence(ctx context.Context, id int64, status string, checkInAt *string) error
}
