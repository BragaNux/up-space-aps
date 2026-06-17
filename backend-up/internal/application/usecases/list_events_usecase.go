package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// EventWithRSVP e um evento junto com a info se o usuario logado ja confirmou presenca
type EventWithRSVP struct {
	*entities.Event
	Confirmed bool `json:"confirmed"`
}

// ListEventsUseCase lista os eventos da agenda
type ListEventsUseCase struct {
	repo ports.EventPort
}

func NewListEventsUseCase(repo ports.EventPort) *ListEventsUseCase {
	return &ListEventsUseCase{repo: repo}
}

// lista os eventos e, se tiver usuario logado, marca quais ele ja confirmou presenca
func (u *ListEventsUseCase) Execute(ctx context.Context, userID *int64) ([]*EventWithRSVP, error) {
	events, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var confirmedIDs map[int64]bool
	if userID != nil {
		confirmedIDs, err = u.repo.RSVPedEventIDs(ctx, *userID)
		if err != nil {
			return nil, err
		}
	}

	result := make([]*EventWithRSVP, 0, len(events))
	for _, event := range events {
		result = append(result, &EventWithRSVP{Event: event, Confirmed: confirmedIDs[event.ID]})
	}
	return result, nil
}
