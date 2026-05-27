package application

import (
	"context"

	"your-next-game-backend/internal/library/domain"
)

type GameRepository interface {
	ListByUserID(ctx context.Context, userID string) ([]domain.Game, error)
	GetByUserIDAndGameID(ctx context.Context, userID, gameID string) (*domain.Game, error)
	CreateGame(ctx context.Context, userID string, game domain.Game) error
}
