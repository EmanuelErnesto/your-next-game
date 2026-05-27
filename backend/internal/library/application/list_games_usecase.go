package application

import (
	"context"

	"your-next-game-backend/internal/library/domain"
)

type ListGamesUseCase struct {
	repo GameRepository
}

func NewListGamesUseCase(repo GameRepository) *ListGamesUseCase {
	return &ListGamesUseCase{repo: repo}
}

func (u *ListGamesUseCase) Execute(ctx context.Context, userID string) ([]domain.Game, error) {
	return u.repo.ListByUserID(ctx, userID)
}
