package usecases

import (
	"context"
	"errors"
	"time"

	"up-espaco/backend/internal/adapters/security"
	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// erro generico de credenciais ruins (de proposito nao diz se foi o email ou a senha que errou)
var ErrInvalidCredentials = errors.New("e-mail ou senha inválidos")

// LoginUseCase confere email/senha e gera o token JWT de acesso
type LoginUseCase struct {
	repo      repositories.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewLoginUseCase(repo repositories.UserRepository, jwtSecret string, tokenTTL time.Duration) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwtSecret: jwtSecret, tokenTTL: tokenTTL}
}

// busca o usuario pelo email, confere a senha com bcrypt e devolve o token JWT
func (u *LoginUseCase) Execute(ctx context.Context, email, password string) (string, *entities.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !security.CheckPassword(user.PasswordHash, password) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := security.GenerateToken(u.jwtSecret, user.ID, user.Role, u.tokenTTL)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
