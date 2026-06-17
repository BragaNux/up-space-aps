package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// PostRepository e a implementacao em postgres do repositorio de posts do feed
type PostRepository struct {
	db *DB
}

func NewPostRepository(db *DB) *PostRepository {
	return &PostRepository{db: db}
}

// List busca os posts de um aluno, incluindo os posts "de turma" que tambem aparecem pros colegas
func (r *PostRepository) List(ctx context.Context, studentID int64) ([]*entities.Post, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `
		SELECT p.id, p.student_id, p.title, p.description, p.pedagogical_note, p.image_url, p.likes, p.bookmarks, p.visibility, s.name as student_name, COALESCE(u.name, s.teacher_name) as teacher_name, COALESCE(u.avatar_url, '') as teacher_avatar_url, (SELECT COUNT(*) FROM comments WHERE comments.post_id = p.id) as comment_count, p.created_at, p.updated_at
		FROM posts p
		JOIN students s ON s.id = p.student_id
		LEFT JOIN users u ON u.id = s.teacher_user_id
		WHERE p.student_id = $1
		   OR (p.visibility = 'turma' AND s.turma_id = (SELECT turma_id FROM students WHERE id = $1))
		ORDER BY p.created_at DESC
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*entities.Post, 0)
	for rows.Next() {
		post := &entities.Post{}
		if err := rows.Scan(&post.ID, &post.StudentID, &post.Title, &post.Description, &post.PedagogicalNote, &post.ImageURL, &post.Likes, &post.Bookmarks, &post.Visibility, &post.StudentName, &post.TeacherName, &post.TeacherAvatar, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

// GetByID busca um post pelo id, ja com nome do aluno e do professor
func (r *PostRepository) GetByID(ctx context.Context, id int64) (*entities.Post, error) {
	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT p.id, p.student_id, p.title, p.description, p.pedagogical_note, p.image_url, p.likes, p.bookmarks, p.visibility, s.name as student_name, COALESCE(u.name, s.teacher_name) as teacher_name, COALESCE(u.avatar_url, '') as teacher_avatar_url, (SELECT COUNT(*) FROM comments WHERE comments.post_id = p.id) as comment_count, p.created_at, p.updated_at
		FROM posts p
		JOIN students s ON s.id = p.student_id
		LEFT JOIN users u ON u.id = s.teacher_user_id
		WHERE p.id = $1
	`, id)

	post := &entities.Post{}
	if err := row.Scan(&post.ID, &post.StudentID, &post.Title, &post.Description, &post.PedagogicalNote, &post.ImageURL, &post.Likes, &post.Bookmarks, &post.Visibility, &post.StudentName, &post.TeacherName, &post.TeacherAvatar, &post.CommentCount, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return nil, err
	}
	return post, nil
}

// Create insere um post novo, com likes/bookmarks zerados e visibilidade "private" como padrao
func (r *PostRepository) Create(ctx context.Context, post *entities.Post) error {
	if post.Visibility == "" {
		post.Visibility = "private"
	}
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO posts (student_id, title, description, pedagogical_note, image_url, likes, bookmarks, visibility, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 0, 0, $6, now(), now())
		RETURNING id, created_at, updated_at
	`, post.StudentID, post.Title, post.Description, post.PedagogicalNote, post.ImageURL, post.Visibility)

	return row.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

// Update atualiza um post existente
func (r *PostRepository) Update(ctx context.Context, post *entities.Post) error {
	if post.Visibility == "" {
		post.Visibility = "private"
	}
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE posts
		SET title = $1, description = $2, pedagogical_note = $3, image_url = $4, visibility = $5, updated_at = now()
		WHERE id = $6
	`, post.Title, post.Description, post.PedagogicalNote, post.ImageURL, post.Visibility, post.ID)
	return err
}

// Delete apaga o post pelo id
func (r *PostRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM posts WHERE id = $1`, id)
	return err
}

// IncrementLikes soma uma curtida e devolve o total atualizado
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

// IncrementBookmarks soma um "salvo" e devolve o total atualizado
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

// DecrementLikes tira uma curtida (nunca deixa ir negativo) e devolve o total atualizado
func (r *PostRepository) DecrementLikes(ctx context.Context, id int64) (int64, error) {
	var likes int64
	row := r.db.Conn().QueryRowContext(ctx, `
		UPDATE posts
		SET likes = GREATEST(0, likes - 1), updated_at = now()
		WHERE id = $1
		RETURNING likes
	`, id)
	if err := row.Scan(&likes); err != nil {
		return 0, err
	}
	return likes, nil
}

// DecrementBookmarks tira um "salvo" (nunca deixa ir negativo) e devolve o total atualizado
func (r *PostRepository) DecrementBookmarks(ctx context.Context, id int64) (int64, error) {
	var bookmarks int64
	row := r.db.Conn().QueryRowContext(ctx, `
		UPDATE posts
		SET bookmarks = GREATEST(0, bookmarks - 1), updated_at = now()
		WHERE id = $1
		RETURNING bookmarks
	`, id)
	if err := row.Scan(&bookmarks); err != nil {
		return 0, err
	}
	return bookmarks, nil
}

// os metodos abaixo so existem pra satisfazer a interface ports.PostPort (que usa nomes diferentes dos do repositories.PostRepository)

func (r *PostRepository) ListPosts(ctx context.Context, studentID int64) ([]*entities.Post, error) {
	return r.List(ctx, studentID)
}

func (r *PostRepository) CreatePost(ctx context.Context, post *entities.Post) error {
	return r.Create(ctx, post)
}

func (r *PostRepository) LikePost(ctx context.Context, id int64) (int64, error) {
	return r.IncrementLikes(ctx, id)
}

func (r *PostRepository) UnlikePost(ctx context.Context, id int64) (int64, error) {
	return r.DecrementLikes(ctx, id)
}

func (r *PostRepository) BookmarkPost(ctx context.Context, id int64) (int64, error) {
	return r.IncrementBookmarks(ctx, id)
}

func (r *PostRepository) UnbookmarkPost(ctx context.Context, id int64) (int64, error) {
	return r.DecrementBookmarks(ctx, id)
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *entities.Post) error {
	return r.Update(ctx, post)
}

func (r *PostRepository) DeletePost(ctx context.Context, id int64) error {
	return r.Delete(ctx, id)
}

var _ repositories.PostRepository = (*PostRepository)(nil)
