package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// TimelineEventRepository cuida dos eventos na linha do tempo de cada aluno
type TimelineEventRepository interface {
	ListForDay(ctx context.Context, studentID int64, date string) ([]*entities.TimelineEvent, error) // lista os eventos do aluno num dia especifico
	GetByID(ctx context.Context, id int64) (*entities.TimelineEvent, error)                          // busca um evento da timeline especifico
	Create(ctx context.Context, event *entities.TimelineEvent) error                                  // cria um evento novo na timeline
	Update(ctx context.Context, event *entities.TimelineEvent) error                                  // atualiza um evento existente
	Delete(ctx context.Context, id int64) error                                                        // remove um evento da timeline
}
