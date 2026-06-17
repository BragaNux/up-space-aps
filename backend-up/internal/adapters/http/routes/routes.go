package routes

import (
	"net/http"
	"time"

	"up-espaco/backend/internal/adapters/http/handlers"
	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/adapters/repositories"
	"up-espaco/backend/internal/application/usecases"

	"github.com/gorilla/mux"
)

// NewRouter monta todas as rotas da API: instancia handlers/usecases e registra cada endpoint no mux
func NewRouter(container *repositories.RepositoryContainer, jwtSecret string, tokenTTL time.Duration) http.Handler {
	requireAuth := middleware.RequireAuth(jwtSecret)
	optionalAuth := middleware.OptionalAuth(jwtSecret)

	authHandler := handlers.NewAuthHandler(
		usecases.NewRegisterUserUseCase(container.UserRepo),
		usecases.NewLoginUseCase(container.UserRepo, jwtSecret, tokenTTL),
		usecases.NewForgotPasswordUseCase(container.UserRepo),
		usecases.NewGetMeUseCase(container.UserRepo),
		usecases.NewUpdateProfileUseCase(container.UserRepo),
		usecases.NewListUsersByRoleUseCase(container.UserRepo),
	)

	studentHandler := handlers.NewStudentHandler(
		usecases.NewGetStudentUseCase(container.StudentRepo),
		usecases.NewGetStudentByIDUseCase(container.StudentRepo),
		usecases.NewListStudentsUseCase(container.StudentRepo),
		usecases.NewListMyChildrenUseCase(container.StudentRepo),
		usecases.NewCreateStudentUseCase(container.StudentRepo),
		usecases.NewUpdateStudentUseCase(container.StudentRepo),
		usecases.NewDeleteStudentUseCase(container.StudentRepo),
		usecases.NewUpdatePresenceUseCase(container.StudentRepo),
		usecases.NewListGuardiansUseCase(container.GuardianRepo),
		container.UserRepo,
		container.TurmaRepo,
	)

	turmaHandler := handlers.NewTurmaHandler(
		usecases.NewListTurmasUseCase(container.TurmaRepo),
		usecases.NewCreateTurmaUseCase(container.TurmaRepo),
		usecases.NewUpdateTurmaUseCase(container.TurmaRepo),
		usecases.NewDeleteTurmaUseCase(container.TurmaRepo),
	)

	guardianHandler := handlers.NewGuardianHandler(
		usecases.NewListGuardiansUseCase(container.GuardianRepo),
		usecases.NewCreateGuardianUseCase(container.GuardianRepo),
		usecases.NewUpdateGuardianUseCase(container.GuardianRepo),
		usecases.NewDeleteGuardianUseCase(container.GuardianRepo),
	)

	timelineHandler := handlers.NewTimelineHandler(
		usecases.NewListTimelineUseCase(container.TimelineRepo),
		usecases.NewGetTimelineEventUseCase(container.TimelineRepo),
		usecases.NewCreateTimelineUseCase(container.TimelineRepo),
		usecases.NewUpdateTimelineUseCase(container.TimelineRepo),
		usecases.NewDeleteTimelineUseCase(container.TimelineRepo),
	)

	postsHandler := handlers.NewPostsHandler(
		usecases.NewListPostsUseCase(container.PostRepo),
		usecases.NewGetPostUseCase(container.PostRepo),
		usecases.NewCreatePostUseCase(container.PostRepo),
		usecases.NewUpdatePostUseCase(container.PostRepo),
		usecases.NewDeletePostUseCase(container.PostRepo),
		usecases.NewLikePostUseCase(container.PostRepo),
		usecases.NewUnlikePostUseCase(container.PostRepo),
		usecases.NewBookmarkPostUseCase(container.PostRepo),
		usecases.NewUnbookmarkPostUseCase(container.PostRepo),
		usecases.NewListCommentsUseCase(container.CommentRepo),
		usecases.NewCreateCommentUseCase(container.CommentRepo),
		usecases.NewDeleteCommentUseCase(container.CommentRepo),
		container.UserRepo,
	)

	eventHandler := handlers.NewEventHandler(
		usecases.NewListEventsUseCase(container.EventRepo),
		usecases.NewGetEventUseCase(container.EventRepo),
		usecases.NewCreateEventUseCase(container.EventRepo),
		usecases.NewUpdateEventUseCase(container.EventRepo),
		usecases.NewDeleteEventUseCase(container.EventRepo),
		usecases.NewRSVPEventUseCase(container.EventRepo),
	)

	announcementHandler := handlers.NewAnnouncementHandler(
		usecases.NewListAnnouncementsUseCase(container.AnnouncementRepo),
		usecases.NewGetAnnouncementUseCase(container.AnnouncementRepo),
		usecases.NewCreateAnnouncementUseCase(container.AnnouncementRepo),
		usecases.NewUpdateAnnouncementUseCase(container.AnnouncementRepo),
		usecases.NewDeleteAnnouncementUseCase(container.AnnouncementRepo),
		usecases.NewMarkAnnouncementReadUseCase(container.AnnouncementRepo),
	)

	attendanceHandler := handlers.NewAttendanceHandler(
		usecases.NewSaveAttendanceUseCase(container.AttendanceRepo),
		usecases.NewGetTurmaAttendanceUseCase(container.AttendanceRepo),
		usecases.NewListAttendanceByStudentUseCase(container.AttendanceRepo),
	)

	milestoneHandler := handlers.NewMilestoneHandler(
		usecases.NewListMilestonesUseCase(container.MilestoneRepo),
		usecases.NewCreateMilestoneUseCase(container.MilestoneRepo),
		usecases.NewUpdateMilestoneUseCase(container.MilestoneRepo),
		usecases.NewDeleteMilestoneUseCase(container.MilestoneRepo),
	)

	router := mux.NewRouter()

	// Auth
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/forgot-password", authHandler.ForgotPassword).Methods(http.MethodPost)
	router.Handle("/api/me", requireAuth(http.HandlerFunc(authHandler.Me))).Methods(http.MethodGet)
	router.Handle("/api/me", requireAuth(http.HandlerFunc(authHandler.UpdateMe))).Methods(http.MethodPut)
	router.Handle("/api/me/children", requireAuth(http.HandlerFunc(studentHandler.ListMyChildren))).Methods(http.MethodGet)
	router.Handle("/api/users", requireAuth(http.HandlerFunc(authHandler.ListByRole))).Methods(http.MethodGet)

	// Students
	router.HandleFunc("/api/student", studentHandler.GetStudent).Methods(http.MethodGet)
	router.HandleFunc("/api/students", studentHandler.List).Methods(http.MethodGet)
	router.Handle("/api/students", requireAuth(http.HandlerFunc(studentHandler.Create))).Methods(http.MethodPost)
	router.HandleFunc("/api/students/{id:[0-9]+}", studentHandler.Get).Methods(http.MethodGet)
	router.Handle("/api/students/{id:[0-9]+}", requireAuth(http.HandlerFunc(studentHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/students/{id:[0-9]+}", requireAuth(http.HandlerFunc(studentHandler.Delete))).Methods(http.MethodDelete)
	router.Handle("/api/students/{id:[0-9]+}/presence", requireAuth(http.HandlerFunc(studentHandler.UpdatePresence))).Methods(http.MethodPatch)
	router.HandleFunc("/api/students/{id:[0-9]+}/guardians", guardianHandler.List).Methods(http.MethodGet)
	router.Handle("/api/students/{id:[0-9]+}/guardians", requireAuth(http.HandlerFunc(guardianHandler.Create))).Methods(http.MethodPost)
	router.Handle("/api/guardians/{id:[0-9]+}", requireAuth(http.HandlerFunc(guardianHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/guardians/{id:[0-9]+}", requireAuth(http.HandlerFunc(guardianHandler.Delete))).Methods(http.MethodDelete)

	// Timeline
	router.HandleFunc("/api/timeline", timelineHandler.List).Methods(http.MethodGet)
	router.Handle("/api/timeline", requireAuth(http.HandlerFunc(timelineHandler.Create))).Methods(http.MethodPost)
	router.HandleFunc("/api/timeline/{id:[0-9]+}", timelineHandler.Get).Methods(http.MethodGet)
	router.Handle("/api/timeline/{id:[0-9]+}", requireAuth(http.HandlerFunc(timelineHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/timeline/{id:[0-9]+}", requireAuth(http.HandlerFunc(timelineHandler.Delete))).Methods(http.MethodDelete)

	// Posts + comments
	router.HandleFunc("/api/posts", postsHandler.List).Methods(http.MethodGet)
	router.Handle("/api/posts", requireAuth(http.HandlerFunc(postsHandler.Create))).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}", postsHandler.Get).Methods(http.MethodGet)
	router.Handle("/api/posts/{id:[0-9]+}", requireAuth(http.HandlerFunc(postsHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/posts/{id:[0-9]+}", requireAuth(http.HandlerFunc(postsHandler.Delete))).Methods(http.MethodDelete)
	router.HandleFunc("/api/posts/{id:[0-9]+}/like", postsHandler.Like).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}/unlike", postsHandler.Unlike).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}/bookmark", postsHandler.Bookmark).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}/unbookmark", postsHandler.Unbookmark).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}/comments", postsHandler.ListComments).Methods(http.MethodGet)
	router.Handle("/api/posts/{id:[0-9]+}/comments", requireAuth(http.HandlerFunc(postsHandler.CreateComment))).Methods(http.MethodPost)
	router.Handle("/api/comments/{id:[0-9]+}", requireAuth(http.HandlerFunc(postsHandler.DeleteComment))).Methods(http.MethodDelete)

	// Events
	router.Handle("/api/events", optionalAuth(http.HandlerFunc(eventHandler.List))).Methods(http.MethodGet)
	router.Handle("/api/events", requireAuth(http.HandlerFunc(eventHandler.Create))).Methods(http.MethodPost)
	router.HandleFunc("/api/events/{id:[0-9]+}", eventHandler.Get).Methods(http.MethodGet)
	router.Handle("/api/events/{id:[0-9]+}", requireAuth(http.HandlerFunc(eventHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/events/{id:[0-9]+}", requireAuth(http.HandlerFunc(eventHandler.Delete))).Methods(http.MethodDelete)
	router.Handle("/api/events/{id:[0-9]+}/rsvp", requireAuth(http.HandlerFunc(eventHandler.RSVP))).Methods(http.MethodPost)

	// Announcements
	router.Handle("/api/announcements", optionalAuth(http.HandlerFunc(announcementHandler.List))).Methods(http.MethodGet)
	router.Handle("/api/announcements", requireAuth(http.HandlerFunc(announcementHandler.Create))).Methods(http.MethodPost)
	router.Handle("/api/announcements/{id:[0-9]+}", optionalAuth(http.HandlerFunc(announcementHandler.Get))).Methods(http.MethodGet)
	router.Handle("/api/announcements/{id:[0-9]+}", requireAuth(http.HandlerFunc(announcementHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/announcements/{id:[0-9]+}", requireAuth(http.HandlerFunc(announcementHandler.Delete))).Methods(http.MethodDelete)
	router.Handle("/api/announcements/{id:[0-9]+}/read", requireAuth(http.HandlerFunc(announcementHandler.MarkRead))).Methods(http.MethodPost)

	// Milestones
	router.HandleFunc("/api/milestones", milestoneHandler.List).Methods(http.MethodGet)
	router.Handle("/api/milestones", requireAuth(http.HandlerFunc(milestoneHandler.Create))).Methods(http.MethodPost)
	router.Handle("/api/milestones/{id:[0-9]+}", requireAuth(http.HandlerFunc(milestoneHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/milestones/{id:[0-9]+}", requireAuth(http.HandlerFunc(milestoneHandler.Delete))).Methods(http.MethodDelete)

	// Turmas
	router.HandleFunc("/api/turmas", turmaHandler.List).Methods(http.MethodGet)
	router.Handle("/api/turmas", requireAuth(http.HandlerFunc(turmaHandler.Create))).Methods(http.MethodPost)
	router.Handle("/api/turmas/{id:[0-9]+}", requireAuth(http.HandlerFunc(turmaHandler.Update))).Methods(http.MethodPut)
	router.Handle("/api/turmas/{id:[0-9]+}", requireAuth(http.HandlerFunc(turmaHandler.Delete))).Methods(http.MethodDelete)

	// Attendance
	router.Handle("/api/students/{id:[0-9]+}/attendance", requireAuth(http.HandlerFunc(attendanceHandler.Save))).Methods(http.MethodPost)
	router.Handle("/api/students/{id:[0-9]+}/attendance", requireAuth(http.HandlerFunc(attendanceHandler.ListByStudent))).Methods(http.MethodGet)
	router.Handle("/api/turmas/{id:[0-9]+}/attendance", requireAuth(http.HandlerFunc(attendanceHandler.ListByTurma))).Methods(http.MethodGet)

	router.Handle("/metrics", middleware.MetricsHandler()).Methods(http.MethodGet)

	// register simple Swagger UI + JSON
	RegisterSwagger(router)

	return router
}
