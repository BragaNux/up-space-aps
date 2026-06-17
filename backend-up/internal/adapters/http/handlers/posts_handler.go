package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// PostsHandler cuida das rotas do feed pedagogico: posts, curtidas, salvos e comentarios
type PostsHandler struct {
	listPosts      *usecases.ListPostsUseCase
	getPost        *usecases.GetPostUseCase
	createPost     *usecases.CreatePostUseCase
	updatePost     *usecases.UpdatePostUseCase
	deletePost     *usecases.DeletePostUseCase
	likePost       *usecases.LikePostUseCase
	unlikePost     *usecases.UnlikePostUseCase
	bookmarkPost   *usecases.BookmarkPostUseCase
	unbookmarkPost *usecases.UnbookmarkPostUseCase
	listComments   *usecases.ListCommentsUseCase
	createComment  *usecases.CreateCommentUseCase
	deleteComment  *usecases.DeleteCommentUseCase
	userRepo       repositories.UserRepository
}

func NewPostsHandler(
	listPosts *usecases.ListPostsUseCase,
	getPost *usecases.GetPostUseCase,
	createPost *usecases.CreatePostUseCase,
	updatePost *usecases.UpdatePostUseCase,
	deletePost *usecases.DeletePostUseCase,
	likePost *usecases.LikePostUseCase,
	unlikePost *usecases.UnlikePostUseCase,
	bookmarkPost *usecases.BookmarkPostUseCase,
	unbookmarkPost *usecases.UnbookmarkPostUseCase,
	listComments *usecases.ListCommentsUseCase,
	createComment *usecases.CreateCommentUseCase,
	deleteComment *usecases.DeleteCommentUseCase,
	userRepo repositories.UserRepository,
) *PostsHandler {
	return &PostsHandler{
		listPosts: listPosts, getPost: getPost, createPost: createPost, updatePost: updatePost,
		deletePost: deletePost, likePost: likePost, unlikePost: unlikePost, bookmarkPost: bookmarkPost, unbookmarkPost: unbookmarkPost,
		listComments: listComments, createComment: createComment, deleteComment: deleteComment,
		userRepo: userRepo,
	}
}

// List devolve os posts do feed, ou lista vazia se nao informar o aluno (GET /api/posts?student_id=)
func (h *PostsHandler) List(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseStudentIDQuery(r)
	if err != nil {
		h.writeJSON(w, http.StatusOK, []any{})
		return
	}

	posts, err := h.listPosts.Execute(r.Context(), studentID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, posts)
}

// Get busca um post especifico (GET /api/posts/{id})
func (h *PostsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	post, err := h.getPost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Post não encontrado")
		return
	}

	h.writeJSON(w, http.StatusOK, post)
}

// Create cria um post novo pra um aluno (POST /api/posts?student_id=)
func (h *PostsHandler) Create(w http.ResponseWriter, r *http.Request) {
	studentID, err := parseStudentIDQuery(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Parâmetro student_id ausente ou inválido")
		return
	}

	var req dto.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	post := &entities.Post{
		StudentID:       studentID,
		Title:           req.Title,
		Description:     req.Description,
		PedagogicalNote: req.PedagogicalNote,
		ImageURL:        req.ImageURL,
		Visibility:      req.Visibility,
	}

	if err := h.createPost.Execute(r.Context(), post); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, post)
}

// Update edita um post existente (PUT /api/posts/{id})
func (h *PostsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	var req dto.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	post := &entities.Post{
		ID:              id,
		Title:           req.Title,
		Description:     req.Description,
		PedagogicalNote: req.PedagogicalNote,
		ImageURL:        req.ImageURL,
		Visibility:      req.Visibility,
	}
	if err := h.updatePost.Execute(r.Context(), post); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, post)
}

// Delete remove um post (DELETE /api/posts/{id})
func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	if err := h.deletePost.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Like da uma curtida num post (POST /api/posts/{id}/like)
func (h *PostsHandler) Like(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	likes, err := h.likePost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"likes": likes})
}

// Bookmark salva um post nos favoritos (POST /api/posts/{id}/bookmark)
func (h *PostsHandler) Bookmark(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	bookmarks, err := h.bookmarkPost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"bookmarks": bookmarks})
}

// ListComments devolve os comentarios de um post (GET /api/posts/{id}/comments)
func (h *PostsHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	postID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	comments, err := h.listComments.Execute(r.Context(), postID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, comments)
}

// CreateComment cria um comentario num post, puxando nome/avatar do usuario logado (POST /api/posts/{id}/comments)
func (h *PostsHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	var req dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Autenticação necessária")
		return
	}

	comment := &entities.Comment{PostID: postID, UserID: &userID, Text: req.Text}
	if user, err := h.userRepo.GetByID(r.Context(), userID); err == nil {
		comment.AuthorName = user.Name
		comment.AvatarURL = user.AvatarURL
	}

	if err := h.createComment.Execute(r.Context(), comment); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, comment)
}

// DeleteComment remove um comentario (DELETE /api/comments/{id})
func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do comentário inválido")
		return
	}

	if err := h.deleteComment.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Unlike tira a curtida de um post (POST /api/posts/{id}/unlike)
func (h *PostsHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	likes, err := h.unlikePost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"likes": likes})
}

// Unbookmark tira um post dos favoritos (POST /api/posts/{id}/unbookmark)
func (h *PostsHandler) Unbookmark(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do post inválido")
		return
	}

	bookmarks, err := h.unbookmarkPost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"bookmarks": bookmarks})
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *PostsHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *PostsHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
