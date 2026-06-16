package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
)

type RSVPEventUseCase struct {
	repo ports.EventPort
}

func NewRSVPEventUseCase(repo ports.EventPort) *RSVPEventUseCase {
	return &RSVPEventUseCase{repo: repo}
}

func (u *RSVPEventUseCase) Execute(ctx context.Context, id int64) (int64, error) {
	return u.repo.RSVPEvent(ctx, id)
}
