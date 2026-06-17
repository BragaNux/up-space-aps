package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// ListAnnouncementsUseCase lista os comunicados disponiveis
type ListAnnouncementsUseCase struct {
	repo repositories.AnnouncementRepository
}

// monta o usecase com o repo de comunicados
func NewListAnnouncementsUseCase(repo repositories.AnnouncementRepository) *ListAnnouncementsUseCase {
	return &ListAnnouncementsUseCase{repo: repo}
}

// pega a lista de comunicados, marcando os ja lidos pelo usuario
func (u *ListAnnouncementsUseCase) Execute(ctx context.Context, userID *int64) ([]*entities.Announcement, error) {
	return u.repo.List(ctx, userID)
}

// GetAnnouncementUseCase busca um comunicado especifico
type GetAnnouncementUseCase struct {
	repo repositories.AnnouncementRepository
}

// monta o usecase com o repo de comunicados
func NewGetAnnouncementUseCase(repo repositories.AnnouncementRepository) *GetAnnouncementUseCase {
	return &GetAnnouncementUseCase{repo: repo}
}

// busca um comunicado pelo id
func (u *GetAnnouncementUseCase) Execute(ctx context.Context, id int64, userID *int64) (*entities.Announcement, error) {
	return u.repo.GetByID(ctx, id, userID)
}

// CreateAnnouncementUseCase cria um comunicado novo
type CreateAnnouncementUseCase struct {
	repo repositories.AnnouncementRepository
}

// monta o usecase com o repo de comunicados
func NewCreateAnnouncementUseCase(repo repositories.AnnouncementRepository) *CreateAnnouncementUseCase {
	return &CreateAnnouncementUseCase{repo: repo}
}

// valida titulo, corpo e prioridade antes de criar o comunicado
func (u *CreateAnnouncementUseCase) Execute(ctx context.Context, a *entities.Announcement) error {
	if a.Title == "" || a.Body == "" {
		return errors.New("título e corpo são obrigatórios")
	}
	if a.Priority != "Urgente" && a.Priority != "Importante" && a.Priority != "Informativo" {
		return errors.New("prioridade inválida")
	}
	return u.repo.Create(ctx, a)
}

// UpdateAnnouncementUseCase atualiza um comunicado existente
type UpdateAnnouncementUseCase struct {
	repo repositories.AnnouncementRepository
}

// monta o usecase com o repo de comunicados
func NewUpdateAnnouncementUseCase(repo repositories.AnnouncementRepository) *UpdateAnnouncementUseCase {
	return &UpdateAnnouncementUseCase{repo: repo}
}

// valida titulo e corpo antes de salvar a atualizacao
func (u *UpdateAnnouncementUseCase) Execute(ctx context.Context, a *entities.Announcement) error {
	if a.Title == "" || a.Body == "" {
		return errors.New("título e corpo são obrigatórios")
	}
	return u.repo.Update(ctx, a)
}

// DeleteAnnouncementUseCase remove um comunicado
type DeleteAnnouncementUseCase struct {
	repo repositories.AnnouncementRepository
}

// monta o usecase com o repo de comunicados
func NewDeleteAnnouncementUseCase(repo repositories.AnnouncementRepository) *DeleteAnnouncementUseCase {
	return &DeleteAnnouncementUseCase{repo: repo}
}

// apaga o comunicado pelo id
func (u *DeleteAnnouncementUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

// MarkAnnouncementReadUseCase marca um comunicado como lido por um usuario
type MarkAnnouncementReadUseCase struct {
	repo repositories.AnnouncementRepository
}

// monta o usecase com o repo de comunicados
func NewMarkAnnouncementReadUseCase(repo repositories.AnnouncementRepository) *MarkAnnouncementReadUseCase {
	return &MarkAnnouncementReadUseCase{repo: repo}
}

// marca que o usuario leu aquele comunicado
func (u *MarkAnnouncementReadUseCase) Execute(ctx context.Context, announcementID, userID int64) error {
	return u.repo.MarkRead(ctx, announcementID, userID)
}
