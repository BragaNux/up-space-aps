package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// UpdatePostUseCase atualiza um post existente do feed
type UpdatePostUseCase struct {
	repo ports.PostPort
}

func NewUpdatePostUseCase(repo ports.PostPort) *UpdatePostUseCase {
	return &UpdatePostUseCase{repo: repo}
}

// confere os campos obrigatorios e a imagem antes de salvar a atualizacao
func (u *UpdatePostUseCase) Execute(ctx context.Context, post *entities.Post) error {
	if post.Title == "" {
		return errors.New("título é obrigatório")
	}
	if post.Description == "" {
		return errors.New("descrição é obrigatória")
	}
	if post.PedagogicalNote == "" {
		return errors.New("nota pedagógica é obrigatória")
	}
	if err := validateImage(post.ImageURL); err != nil {
		return err
	}
	return u.repo.UpdatePost(ctx, post)
}
