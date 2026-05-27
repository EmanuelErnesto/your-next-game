package application

import (
	"context"

	"github.com/google/uuid"
	"your-next-game-backend/internal/identity/domain"
)

type CreateUserUseCase struct {
	repo domain.UserRepository
}

func NewCreateUserUseCase(repo domain.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo: repo,
	}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, user domain.User) (*domain.User, error) {
	existing, err := u.repo.GetBySteamID(ctx, user.SteamID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Update profile info if changed
		if existing.Name != user.Name || existing.AvatarFull != user.AvatarFull {
			existing.Name = user.Name
			existing.AvatarFull = user.AvatarFull
			if err := u.repo.Save(ctx, *existing); err != nil {
				return nil, err
			}
		}
		return existing, nil
	}

	// Generate UUID for the new user
	user.ID = uuid.New().String()
	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return &user, nil
}

