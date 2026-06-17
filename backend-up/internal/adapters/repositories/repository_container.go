package repositories

import (
	"up-espaco/backend/internal/adapters/database/postgres"
)

// RepositoryContainer junta todos os repositorios postgres num lugar so, pra facilitar a injecao de dependencia no main
type RepositoryContainer struct {
	StudentRepo      *postgres.StudentRepository
	TimelineRepo     *postgres.TimelineEventRepository
	PostRepo         *postgres.PostRepository
	EventRepo        *postgres.EventRepository
	UserRepo         *postgres.UserRepository
	CommentRepo      *postgres.CommentRepository
	AnnouncementRepo *postgres.AnnouncementRepository
	GuardianRepo     *postgres.GuardianRepository
	MilestoneRepo    *postgres.MilestoneRepository
	TurmaRepo        *postgres.TurmaRepository
	AttendanceRepo   *postgres.AttendanceRepository
}

// NewRepositoryContainer instancia todos os repositorios em cima da mesma conexao de banco
func NewRepositoryContainer(db *postgres.DB) *RepositoryContainer {
	return &RepositoryContainer{
		StudentRepo:      postgres.NewStudentRepository(db),
		TimelineRepo:     postgres.NewTimelineEventRepository(db),
		PostRepo:         postgres.NewPostRepository(db),
		EventRepo:        postgres.NewEventRepository(db),
		UserRepo:         postgres.NewUserRepository(db),
		CommentRepo:      postgres.NewCommentRepository(db),
		AnnouncementRepo: postgres.NewAnnouncementRepository(db),
		GuardianRepo:     postgres.NewGuardianRepository(db),
		MilestoneRepo:    postgres.NewMilestoneRepository(db),
		TurmaRepo:        postgres.NewTurmaRepository(db),
		AttendanceRepo:   postgres.NewAttendanceRepository(db),
	}
}
