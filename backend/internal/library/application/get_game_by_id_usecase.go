package application

import (
	"context"

	"your-next-game-backend/internal/library/domain"
)

type GetGameByIDUseCase struct {
	repo GameRepository
}

func NewGetGameByIDUseCase(repo GameRepository) *GetGameByIDUseCase {
	return &GetGameByIDUseCase{repo: repo}
}

func (u *GetGameByIDUseCase) Execute(
	ctx context.Context,
	userID string,
	gameID string,
) (*domain.Game, error) {
	return u.repo.GetByUserIDAndGameID(ctx, userID, gameID)
}
