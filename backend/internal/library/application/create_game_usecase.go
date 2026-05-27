package application

import (
	"context"

	"your-next-game-backend/internal/library/domain"
)

type CreateGameUseCase struct {
	repo GameRepository
}

func NewCreateGameUseCase(repo GameRepository) *CreateGameUseCase {
	return &CreateGameUseCase{
		repo: repo,
	}
}

func (u *CreateGameUseCase) Execute(ctx context.Context, userID string, game domain.Game) error {
	// Aqui poderíamos ter mais lógicas de negócio, como validar limites, 
	// disparar eventos, etc.
	return u.repo.CreateGame(ctx, userID, game)
}
