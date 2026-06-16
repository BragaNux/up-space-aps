package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/application/ports"
	"up-espaco/backend/internal/domain/entities"
)

type CreatePostUseCase struct {
	repo ports.PostPort
}

func NewCreatePostUseCase(repo ports.PostPort) *CreatePostUseCase {
	return &CreatePostUseCase{repo: repo}
}

func (u *CreatePostUseCase) Execute(ctx context.Context, post *entities.Post) error {
	if post.Title == "" {
		return errors.New("title is required")
	}
	if post.Description == "" {
		return errors.New("description is required")
	}
	if post.PedagogicalNote == "" {
		return errors.New("pedagogical_note is required")
	}

	return u.repo.CreatePost(ctx, post)
}
