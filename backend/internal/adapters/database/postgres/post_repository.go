package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

type PostRepository struct {
	db *DB
}

func NewPostRepository(db *DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) List(ctx context.Context) ([]*entities.Post, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT id, title, description, pedagogical_note, likes, bookmarks, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*entities.Post, 0)
	for rows.Next() {
		post := &entities.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.PedagogicalNote, &post.Likes, &post.Bookmarks, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (r *PostRepository) Create(ctx context.Context, post *entities.Post) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO posts (title, description, pedagogical_note, likes, bookmarks, created_at, updated_at)
		VALUES ($1, $2, $3, 0, 0, now(), now())
		RETURNING id, created_at, updated_at
	`, post.Title, post.Description, post.PedagogicalNote)

	return row.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (r *PostRepository) IncrementLikes(ctx context.Context, id int64) (int64, error) {
	var likes int64
	row := r.db.Conn().QueryRowContext(ctx, `
		UPDATE posts
		SET likes = likes + 1, updated_at = now()
		WHERE id = $1
		RETURNING likes
	`, id)
	if err := row.Scan(&likes); err != nil {
		return 0, err
	}
	return likes, nil
}

func (r *PostRepository) IncrementBookmarks(ctx context.Context, id int64) (int64, error) {
	var bookmarks int64
	row := r.db.Conn().QueryRowContext(ctx, `
		UPDATE posts
		SET bookmarks = bookmarks + 1, updated_at = now()
		WHERE id = $1
		RETURNING bookmarks
	`, id)
	if err := row.Scan(&bookmarks); err != nil {
		return 0, err
	}
	return bookmarks, nil
}

func (r *PostRepository) ListPosts(ctx context.Context) ([]*entities.Post, error) {
	return r.List(ctx)
}

func (r *PostRepository) CreatePost(ctx context.Context, post *entities.Post) error {
	return r.Create(ctx, post)
}

func (r *PostRepository) LikePost(ctx context.Context, id int64) (int64, error) {
	return r.IncrementLikes(ctx, id)
}

func (r *PostRepository) BookmarkPost(ctx context.Context, id int64) (int64, error) {
	return r.IncrementBookmarks(ctx, id)
}

var _ repositories.PostRepository = (*PostRepository)(nil)
