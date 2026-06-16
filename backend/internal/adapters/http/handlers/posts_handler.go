package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"

	"github.com/gorilla/mux"
)

type PostsHandler struct {
	listPosts    *usecases.ListPostsUseCase
	createPost   *usecases.CreatePostUseCase
	likePost     *usecases.LikePostUseCase
	bookmarkPost *usecases.BookmarkPostUseCase
}

func NewPostsHandler(listPosts *usecases.ListPostsUseCase, createPost *usecases.CreatePostUseCase, likePost *usecases.LikePostUseCase, bookmarkPost *usecases.BookmarkPostUseCase) *PostsHandler {
	return &PostsHandler{listPosts: listPosts, createPost: createPost, likePost: likePost, bookmarkPost: bookmarkPost}
}

func (h *PostsHandler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.listPosts.Execute(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, posts)
}

func (h *PostsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post := &entities.Post{
		Title:           req.Title,
		Description:     req.Description,
		PedagogicalNote: req.PedagogicalNote,
	}

	if err := h.createPost.Execute(r.Context(), post); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, post)
}

func (h *PostsHandler) Like(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		h.writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	likes, err := h.likePost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"likes": likes})
}

func (h *PostsHandler) Bookmark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		h.writeError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	bookmarks, err := h.bookmarkPost.Execute(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]int64{"bookmarks": bookmarks})
}

func (h *PostsHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *PostsHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
