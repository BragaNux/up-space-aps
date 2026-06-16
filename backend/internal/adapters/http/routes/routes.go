package routes

import (
	"net/http"

	"up-espaco/backend/internal/adapters/http/handlers"
	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/adapters/repositories"
	"up-espaco/backend/internal/application/usecases"

	"github.com/gorilla/mux"
)

func NewRouter(container *repositories.RepositoryContainer) http.Handler {
	studentHandler := handlers.NewStudentHandler(
		usecases.NewGetStudentUseCase(container.StudentRepo),
		usecases.NewUpdatePresenceUseCase(container.StudentRepo),
	)

	timelineHandler := handlers.NewTimelineHandler(
		usecases.NewListTimelineUseCase(container.TimelineRepo),
	)

	postsHandler := handlers.NewPostsHandler(
		usecases.NewListPostsUseCase(container.PostRepo),
		usecases.NewCreatePostUseCase(container.PostRepo),
		usecases.NewLikePostUseCase(container.PostRepo),
		usecases.NewBookmarkPostUseCase(container.PostRepo),
	)

	eventHandler := handlers.NewEventHandler(
		usecases.NewListEventsUseCase(container.EventRepo),
		usecases.NewRSVPEventUseCase(container.EventRepo),
	)

	router := mux.NewRouter()

	router.HandleFunc("/api/student", studentHandler.GetStudent).Methods(http.MethodGet)
	router.HandleFunc("/api/student/presence", studentHandler.UpdatePresence).Methods(http.MethodPatch)
	router.HandleFunc("/api/timeline", timelineHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/api/posts", postsHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/api/posts", postsHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}/like", postsHandler.Like).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{id:[0-9]+}/bookmark", postsHandler.Bookmark).Methods(http.MethodPost)
	router.HandleFunc("/api/events", eventHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/api/events/{id:[0-9]+}/rsvp", eventHandler.RSVP).Methods(http.MethodPost)
	router.Handle("/metrics", middleware.MetricsHandler()).Methods(http.MethodGet)

	// register simple Swagger UI + JSON
	RegisterSwagger(router)

	return router
}
