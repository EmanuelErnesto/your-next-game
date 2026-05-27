package application

import (
	"context"
	"testing"

	"your-next-game-backend/internal/library/domain"
)

type fakeRepository struct {
	games []domain.Game
}

func (f *fakeRepository) ListByUserID(_ context.Context, _ string) ([]domain.Game, error) {
	return f.games, nil
}

func (f *fakeRepository) GetByUserIDAndGameID(_ context.Context, _, _ string) (*domain.Game, error) {
	if len(f.games) == 0 {
		return nil, nil
	}
	game := f.games[0]
	return &game, nil
}

func (f *fakeRepository) CreateGame(_ context.Context, _ string, game domain.Game) error {
	f.games = append(f.games, game)
	return nil
}

func TestListGamesUseCaseExecute(t *testing.T) {
	repo := &fakeRepository{
		games: []domain.Game{
			{ID: "1", Title: "Game"},
		},
	}
	useCase := NewListGamesUseCase(repo)

	games, err := useCase.Execute(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(games))
	}
}
