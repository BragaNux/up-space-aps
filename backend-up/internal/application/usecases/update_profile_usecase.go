package usecases

import (
	"context"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// UpdateProfileUseCase atualiza os dados do perfil do usuario logado
type UpdateProfileUseCase struct {
	repo repositories.UserRepository
}

func NewUpdateProfileUseCase(repo repositories.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{repo: repo}
}

// so troca os campos que vieram preenchidos, mantendo o resto como ja estava
func (u *UpdateProfileUseCase) Execute(ctx context.Context, userID int64, name, phone, address, avatarURL string) (*entities.User, error) {
	if err := validateImage(avatarURL); err != nil {
		return nil, err
	}

	user, err := u.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if phone != "" {
		user.Phone = phone
	}
	if address != "" {
		user.Address = address
	}
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}

	if err := u.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
