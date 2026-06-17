package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// EventRepository cuida dos eventos da agenda e dos RSVPs (confirmacoes de presenca)
type EventRepository interface {
	List(ctx context.Context) ([]*entities.Event, error)                                                  // lista todos os eventos
	GetByID(ctx context.Context, id int64) (*entities.Event, error)                                       // busca um evento especifico
	Create(ctx context.Context, event *entities.Event) error                                              // cria um evento novo
	Update(ctx context.Context, event *entities.Event) error                                              // atualiza um evento existente
	Delete(ctx context.Context, id int64) error                                                            // remove um evento
	ToggleRSVP(ctx context.Context, eventID int64, userID int64) (count int64, confirmed bool, err error) // confirma ou cancela presenca do usuario no evento
	RSVPedEventIDs(ctx context.Context, userID int64) (map[int64]bool, error)                              // diz em quais eventos o usuario ja confirmou presenca
}
