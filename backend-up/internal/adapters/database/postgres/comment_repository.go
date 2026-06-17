package postgres

import (
	"context"
	"database/sql"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// CommentRepository e a implementacao em postgres do repositorio de comentarios
type CommentRepository struct {
	db *DB
}

func NewCommentRepository(db *DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// ListByPost busca os comentarios de um post em ordem cronologica
func (r *CommentRepository) ListByPost(ctx context.Context, postID int64) ([]*entities.Comment, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, post_id, user_id, author_name, avatar_url, text, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]*entities.Comment, 0)
	for rows.Next() {
		comment, err := scanComment(rows.Scan)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

// scanComment le uma linha do banco pra dentro de um Comment, tratando user_id nulo (comentario sem usuario vinculado)
func scanComment(scan func(dest ...any) error) (*entities.Comment, error) {
	comment := &entities.Comment{}
	var userID sql.NullInt64
	if err := scan(&comment.ID, &comment.PostID, &userID, &comment.AuthorName, &comment.AvatarURL, &comment.Text, &comment.CreatedAt); err != nil {
		return nil, err
	}
	if userID.Valid {
		id := userID.Int64
		comment.UserID = &id
	}
	return comment, nil
}

// GetByID busca um comentario pelo id
func (r *CommentRepository) GetByID(ctx context.Context, id int64) (*entities.Comment, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, post_id, user_id, author_name, avatar_url, text, created_at
		FROM comments
		WHERE id = $1
	`, id)
	return scanComment(row.Scan)
}

// Create insere um comentario novo
func (r *CommentRepository) Create(ctx context.Context, comment *entities.Comment) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO comments (post_id, user_id, author_name, avatar_url, text)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, comment.PostID, comment.UserID, comment.AuthorName, comment.AvatarURL, comment.Text)
	return row.Scan(&comment.ID, &comment.CreatedAt)
}

// Delete apaga o comentario pelo id
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM comments WHERE id = $1`, id)
	return err
}

var _ repositories.CommentRepository = (*CommentRepository)(nil)
