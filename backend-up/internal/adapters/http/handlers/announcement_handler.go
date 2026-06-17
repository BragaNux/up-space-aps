package handlers

import (
	"encoding/json"
	"net/http"

	"up-espaco/backend/internal/adapters/http/middleware"
	"up-espaco/backend/internal/application/dto"
	"up-espaco/backend/internal/application/usecases"
	"up-espaco/backend/internal/domain/entities"
)

// AnnouncementHandler cuida das rotas de comunicados/avisos
type AnnouncementHandler struct {
	list     *usecases.ListAnnouncementsUseCase
	get      *usecases.GetAnnouncementUseCase
	create   *usecases.CreateAnnouncementUseCase
	update   *usecases.UpdateAnnouncementUseCase
	delete   *usecases.DeleteAnnouncementUseCase
	markRead *usecases.MarkAnnouncementReadUseCase
}

func NewAnnouncementHandler(
	list *usecases.ListAnnouncementsUseCase,
	get *usecases.GetAnnouncementUseCase,
	create *usecases.CreateAnnouncementUseCase,
	update *usecases.UpdateAnnouncementUseCase,
	delete *usecases.DeleteAnnouncementUseCase,
	markRead *usecases.MarkAnnouncementReadUseCase,
) *AnnouncementHandler {
	return &AnnouncementHandler{list: list, get: get, create: create, update: update, delete: delete, markRead: markRead}
}

// userIDPointer devolve o id do usuario logado, ou nil se a rota for anonima
func userIDPointer(r *http.Request) *int64 {
	if id, ok := middleware.UserIDFromContext(r.Context()); ok {
		return &id
	}
	return nil
}

// requireProfissional confere se quem ta chamando e uma conta profissional
func requireProfissional(r *http.Request) bool {
	role, ok := middleware.RoleFromContext(r.Context())
	return ok && role == "profissional"
}

// List devolve os comunicados, marcando os ja lidos pelo usuario logado (GET /api/announcements)
func (h *AnnouncementHandler) List(w http.ResponseWriter, r *http.Request) {
	announcements, err := h.list.Execute(r.Context(), userIDPointer(r))
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, announcements)
}

// Get busca um comunicado especifico (GET /api/announcements/{id})
func (h *AnnouncementHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do comunicado inválido")
		return
	}

	announcement, err := h.get.Execute(r.Context(), id, userIDPointer(r))
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Comunicado não encontrado")
		return
	}

	h.writeJSON(w, http.StatusOK, announcement)
}

// Create cria um comunicado novo, so pra contas profissionais (POST /api/announcements)
func (h *AnnouncementHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem criar comunicados")
		return
	}

	var req dto.CreateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	announcement := &entities.Announcement{
		Title: req.Title, Sender: req.Sender, Priority: req.Priority,
		Preview: req.Preview, Body: req.Body, AttachmentName: req.AttachmentName,
	}

	if err := h.create.Execute(r.Context(), announcement); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusCreated, announcement)
}

// Update edita um comunicado existente, so pra contas profissionais (PUT /api/announcements/{id})
func (h *AnnouncementHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem editar comunicados")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do comunicado inválido")
		return
	}

	var req dto.CreateAnnouncementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Corpo da requisição inválido")
		return
	}

	announcement := &entities.Announcement{
		ID: id, Title: req.Title, Sender: req.Sender, Priority: req.Priority,
		Preview: req.Preview, Body: req.Body, AttachmentName: req.AttachmentName,
	}

	if err := h.update.Execute(r.Context(), announcement); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, announcement)
}

// Delete remove um comunicado, so pra contas profissionais (DELETE /api/announcements/{id})
func (h *AnnouncementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if !requireProfissional(r) {
		h.writeError(w, http.StatusForbidden, "Somente contas profissionais podem remover comunicados")
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do comunicado inválido")
		return
	}

	if err := h.delete.Execute(r.Context(), id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// MarkRead marca o comunicado como lido pelo usuario logado (POST /api/announcements/{id}/read)
func (h *AnnouncementHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "ID do comunicado inválido")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "Autenticação necessária")
		return
	}

	if err := h.markRead.Execute(r.Context(), id, userID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"status": "lido"})
}

// writeJSON serializa qualquer payload como JSON com o status code informado
func (h *AnnouncementHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// writeError monta uma resposta de erro padrao {"error": "..."}
func (h *AnnouncementHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}
