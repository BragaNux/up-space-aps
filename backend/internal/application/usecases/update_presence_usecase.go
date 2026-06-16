package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	domainErrors "up-espaco/backend/internal/domain/errors"
)

type UpdatePresenceUseCase struct {
	repo ports.StudentPort
}

func NewUpdatePresenceUseCase(repo ports.StudentPort) *UpdatePresenceUseCase {
	return &UpdatePresenceUseCase{repo: repo}
}

func (u *UpdatePresenceUseCase) Execute(ctx context.Context, status string) error {
	if status != "present" && status != "absent" {
		return domainErrors.ErrInvalidPresenceStatus
	}

	student, err := u.repo.GetActive(ctx)
	if err != nil {
		return err
	}

	return u.repo.UpdatePresence(ctx, student.ID, status, nil)
}
