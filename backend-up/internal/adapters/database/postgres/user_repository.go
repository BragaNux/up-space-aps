package postgres

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// UserRepository e a implementacao em postgres do repositorio de usuarios
type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

const userColumns = `id, name, email, password_hash, role, phone, address, avatar_url, created_at, updated_at`

// scanUser le uma linha do banco pra dentro de um User
func scanUser(scan func(dest ...any) error) (*entities.User, error) {
	user := &entities.User{}
	if err := scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Phone, &user.Address, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

// Create insere um usuario novo (a senha ja deve chegar hasheada)
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO users (name, email, password_hash, role, phone, address, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`, user.Name, user.Email, user.PasswordHash, user.Role, user.Phone, user.Address, user.AvatarURL)
	return row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetByEmail busca um usuario pelo email, usado no login
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	row := r.db.Conn().QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE email = $1`, email)
	return scanUser(row.Scan)
}

// GetByID busca um usuario pelo id
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entities.User, error) {
	row := r.db.Conn().QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row.Scan)
}

// ListByRole busca usuarios de um papel especifico, ordenados por nome
func (r *UserRepository) ListByRole(ctx context.Context, role string) ([]*entities.User, error) {
	rows, err := r.db.Conn().QueryContext(ctx, `SELECT `+userColumns+` FROM users WHERE role = $1 ORDER BY name ASC`, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*entities.User, 0)
	for rows.Next() {
		user, err := scanUser(rows.Scan)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Update atualiza os dados de perfil do usuario (nao toca em email/senha/role)
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE users
		SET name = $1, phone = $2, address = $3, avatar_url = $4, updated_at = now()
		WHERE id = $5
	`, user.Name, user.Phone, user.Address, user.AvatarURL, user.ID)
	return err
}

var _ repositories.UserRepository = (*UserRepository)(nil)
