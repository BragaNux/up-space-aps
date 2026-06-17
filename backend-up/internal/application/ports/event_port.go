package ports

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// EventPort e a porta que os usecases de evento usam pra falar com o repositorio (evita import direto do pacote repositories)
type EventPort interface {
	List(ctx context.Context) ([]*entities.Event, error)
	GetByID(ctx context.Context, id int64) (*entities.Event, error)
	Create(ctx context.Context, event *entities.Event) error
	Update(ctx context.Context, event *entities.Event) error
	Delete(ctx context.Context, id int64) error
	ToggleRSVP(ctx context.Context, eventID int64, userID int64) (count int64, confirmed bool, err error)
	RSVPedEventIDs(ctx context.Context, userID int64) (map[int64]bool, error)
}
