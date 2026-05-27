package application

import (
	"context"
	"testing"

	"your-next-game-backend/internal/library/domain"
)

func TestGetGameByIDUseCaseExecute(t *testing.T) {
	repo := &fakeRepository{
		games: []domain.Game{
			{ID: "1", Title: "Game"},
		},
	}
	useCase := NewGetGameByIDUseCase(repo)

	game, err := useCase.Execute(context.Background(), "user-1", "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if game == nil {
		t.Fatal("expected game, got nil")
	}
	if game.ID != "1" {
		t.Fatalf("expected game id 1, got %s", game.ID)
	}
}
