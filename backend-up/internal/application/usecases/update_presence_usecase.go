package usecases

import (
	"context"

	"up-espaco/backend/internal/application/ports"
	domainErrors "up-espaco/backend/internal/domain/errors"
)

// UpdatePresenceUseCase marca um aluno como presente ou ausente
type UpdatePresenceUseCase struct {
	repo ports.StudentPort
}

func NewUpdatePresenceUseCase(repo ports.StudentPort) *UpdatePresenceUseCase {
	return &UpdatePresenceUseCase{repo: repo}
}

// confere se o status e valido e atualiza a presenca do aluno
func (u *UpdatePresenceUseCase) Execute(ctx context.Context, studentID int64, status string) error {
	if status != "present" && status != "absent" {
		return domainErrors.ErrInvalidPresenceStatus
	}

	return u.repo.UpdatePresence(ctx, studentID, status, nil)
}
