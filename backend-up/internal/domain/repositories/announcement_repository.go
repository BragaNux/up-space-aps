package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// AnnouncementRepository define como buscar/salvar comunicados no banco
type AnnouncementRepository interface {
	List(ctx context.Context, userID *int64) ([]*entities.Announcement, error)                 // lista os comunicados, marcando quais o usuario ja leu
	GetByID(ctx context.Context, id int64, userID *int64) (*entities.Announcement, error)       // busca um comunicado especifico
	Create(ctx context.Context, announcement *entities.Announcement) error                      // cria um comunicado novo
	Update(ctx context.Context, announcement *entities.Announcement) error                      // atualiza um comunicado existente
	Delete(ctx context.Context, id int64) error                                                 // remove um comunicado
	MarkRead(ctx context.Context, announcementID int64, userID int64) error                     // marca o comunicado como lido pelo usuario
}
