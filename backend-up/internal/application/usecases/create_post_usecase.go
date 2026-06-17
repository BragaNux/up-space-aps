package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

// CreatePostUseCase cria um post novo no feed pedagogico
type CreatePostUseCase struct {
	repo ports.PostPort
}

func NewCreatePostUseCase(repo ports.PostPort) *CreatePostUseCase {
	return &CreatePostUseCase{repo: repo}
}

// confere os campos obrigatorios e a imagem antes de criar o post
func (u *CreatePostUseCase) Execute(ctx context.Context, post *entities.Post) error {
	if post.StudentID == 0 {
		return errors.New("id do aluno é obrigatório")
	}
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

	return u.repo.CreatePost(ctx, post)
}
