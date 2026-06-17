package usecases

import (
	"context"

	"up-espaco/backend/internal/domain/repositories"
)

// ForgotPasswordUseCase recebe o pedido de "esqueci a senha" sem dizer se o email existe ou nao
// (ainda nao tem envio de email implementado, isso aqui e so um placeholder)
type ForgotPasswordUseCase struct {
	repo repositories.UserRepository
}

func NewForgotPasswordUseCase(repo repositories.UserRepository) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{repo: repo}
}

// so confere se o email existe, mas nunca retorna erro pra nao revelar isso pro usuario
func (u *ForgotPasswordUseCase) Execute(ctx context.Context, email string) error {
	_, _ = u.repo.GetByEmail(ctx, email)
	return nil
}
