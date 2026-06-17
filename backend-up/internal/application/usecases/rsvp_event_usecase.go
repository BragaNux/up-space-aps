package usecases

import (
	"context"
	"errors"
	"time"

	"up-espaco/backend/internal/application/ports"
)

// RSVPEventUseCase confirma ou cancela a presenca do usuario num evento
type RSVPEventUseCase struct {
	repo ports.EventPort
}

func NewRSVPEventUseCase(repo ports.EventPort) *RSVPEventUseCase {
	return &RSVPEventUseCase{repo: repo}
}

// nao deixa confirmar presenca em evento que ja passou, e inverte o RSVP do usuario
func (u *RSVPEventUseCase) Execute(ctx context.Context, eventID int64, userID int64) (count int64, confirmed bool, err error) {
	event, err := u.repo.GetByID(ctx, eventID)
	if err != nil {
		return 0, false, err
	}
	if event.EndsAt.Before(time.Now()) {
		return 0, false, errors.New("Não é possível confirmar presença em eventos que já passaram")
	}
	return u.repo.ToggleRSVP(ctx, eventID, userID)
}

