package usecases

import (
	"context"
	"time"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// ListTimelineUseCase lista os eventos do dia de hoje na timeline do aluno
type ListTimelineUseCase struct {
	repo ports.TimelinePort
}

func NewListTimelineUseCase(repo ports.TimelinePort) *ListTimelineUseCase {
	return &ListTimelineUseCase{repo: repo}
}

// busca os eventos da timeline do aluno no dia de hoje
func (u *ListTimelineUseCase) Execute(ctx context.Context, studentID int64) ([]*entities.TimelineEvent, error) {
	return u.repo.ListForDay(ctx, studentID, time.Now().Format("2006-01-02"))
}
