package usecases

import (
	"context"
	"errors"

	"up-espaco/backend/internal/adapters/security"
	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

var ErrEmailAlreadyInUse = errors.New("e-mail já está em uso")

// RegisterUserUseCase cria uma conta nova (cadastro de profissional ou responsavel)
type RegisterUserUseCase struct {
	repo repositories.UserRepository
}

func NewRegisterUserUseCase(repo repositories.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo}
}

// valida os dados, confere se o email ja existe, hasheia a senha com bcrypt e cria o usuario
func (u *RegisterUserUseCase) Execute(ctx context.Context, name, email, password, role string) (*entities.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("nome, e-mail e senha são obrigatórios")
	}
	if role != "profissional" && role != "responsavel" {
		return nil, errors.New("papel inválido")
	}

	if existing, err := u.repo.GetByEmail(ctx, email); err == nil && existing != nil {
		return nil, ErrEmailAlreadyInUse
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{Name: name, Email: email, PasswordHash: hash, Role: role}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
