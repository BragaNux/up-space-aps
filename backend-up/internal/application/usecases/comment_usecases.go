package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// ListCommentsUseCase lista os comentarios de um post
type ListCommentsUseCase struct {
	repo repositories.CommentRepository
}

func NewListCommentsUseCase(repo repositories.CommentRepository) *ListCommentsUseCase {
	return &ListCommentsUseCase{repo: repo}
}

// busca os comentarios daquele post
func (u *ListCommentsUseCase) Execute(ctx context.Context, postID int64) ([]*entities.Comment, error) {
	return u.repo.ListByPost(ctx, postID)
}

// CreateCommentUseCase cria um comentario novo num post
type CreateCommentUseCase struct {
	repo repositories.CommentRepository
}

func NewCreateCommentUseCase(repo repositories.CommentRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{repo: repo}
}

// exige que o comentario tenha texto antes de salvar
func (u *CreateCommentUseCase) Execute(ctx context.Context, comment *entities.Comment) error {
	if comment.Text == "" {
		return errors.New("texto é obrigatório")
	}
	return u.repo.Create(ctx, comment)
}

// DeleteCommentUseCase remove um comentario
type DeleteCommentUseCase struct {
	repo repositories.CommentRepository
}

func NewDeleteCommentUseCase(repo repositories.CommentRepository) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{repo: repo}
}

// apaga o comentario pelo id
func (u *DeleteCommentUseCase) Execute(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
