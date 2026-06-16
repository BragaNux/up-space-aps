package repositories

import (
	"up-espaco/backend/internal/adapters/database/postgres"
)

type RepositoryContainer struct {
	StudentRepo  *postgres.StudentRepository
	TimelineRepo *postgres.TimelineEventRepository
	PostRepo     *postgres.PostRepository
	EventRepo    *postgres.EventRepository
}

func NewRepositoryContainer(db *postgres.DB) *RepositoryContainer {
	return &RepositoryContainer{
		StudentRepo:  postgres.NewStudentRepository(db),
		TimelineRepo: postgres.NewTimelineEventRepository(db),
		PostRepo:     postgres.NewPostRepository(db),
		EventRepo:    postgres.NewEventRepository(db),
	}
}
