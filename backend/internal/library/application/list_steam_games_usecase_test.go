package application

import (
	"context"
	"testing"

	"your-next-game-backend/internal/platform/steamapi"
)

type fakeSteamProvider struct {
	games []steamapi.OwnedGame
}

func (f *fakeSteamProvider) GetOwnedGames(_ context.Context, _ string) ([]steamapi.OwnedGame, error) {
	return f.games, nil
}

func TestListSteamGamesUseCaseExecute(t *testing.T) {
	provider := &fakeSteamProvider{
		games: []steamapi.OwnedGame{
			{AppID: 10, Name: "Counter-Strike", PlaytimeForever: 120},
		},
	}
	useCase := NewListSteamGamesUseCase(provider)

	games, err := useCase.Execute(context.Background(), "steam-user")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(games) != 1 {
		t.Fatalf("expected 1 game, got %d", len(games))
	}
	if games[0].ID != "10" {
		t.Fatalf("expected id 10, got %s", games[0].ID)
	}
}
