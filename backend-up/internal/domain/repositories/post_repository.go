package repositories

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
)

// PostRepository cuida dos posts do feed pedagogico (curtidas e salvos inclusos)
type PostRepository interface {
	List(ctx context.Context, studentID int64) ([]*entities.Post, error) // lista os posts (de um aluno especifico, se informado)
	GetByID(ctx context.Context, id int64) (*entities.Post, error)      // busca um post especifico
	Create(ctx context.Context, post *entities.Post) error              // cria um post novo
	Update(ctx context.Context, post *entities.Post) error              // atualiza um post existente
	Delete(ctx context.Context, id int64) error                         // remove um post
	IncrementLikes(ctx context.Context, id int64) (int64, error)        // soma uma curtida e devolve o total atualizado
	DecrementLikes(ctx context.Context, id int64) (int64, error)        // tira uma curtida e devolve o total atualizado
	IncrementBookmarks(ctx context.Context, id int64) (int64, error)    // soma um "salvo" e devolve o total atualizado
	DecrementBookmarks(ctx context.Context, id int64) (int64, error)    // tira um "salvo" e devolve o total atualizado
}
