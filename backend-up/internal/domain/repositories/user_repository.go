package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// UserRepository cuida do cadastro dos usuarios do sistema (login, perfil etc)
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error                       // cria um usuario novo
	GetByEmail(ctx context.Context, email string) (*entities.User, error)       // busca um usuario pelo email (usado no login)
	GetByID(ctx context.Context, id int64) (*entities.User, error)              // busca um usuario pelo id
	Update(ctx context.Context, user *entities.User) error                       // atualiza dados de um usuario
	ListByRole(ctx context.Context, role string) ([]*entities.User, error)      // lista usuarios filtrando por papel (admin, professor, responsavel)
}
