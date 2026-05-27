package application

import (
	"context"
	"testing"

	"your-next-game-backend/internal/platform/steamapi"
)

type fakeStoreProvider struct {
	details *steamapi.StoreGameDetails
}

func (f *fakeStoreProvider) GetStoreDetails(_ context.Context, _ string) (*steamapi.StoreGameDetails, error) {
	return f.details, nil
}

func TestGetSteamStoreDetailsUseCaseExecute(t *testing.T) {
	provider := &fakeStoreProvider{
		details: &steamapi.StoreGameDetails{
			Name:                "Portal 2",
			DetailedDescription: "Desc",
			Developers:          []string{"Valve"},
			HeaderImage:         "https://img",
		},
	}
	useCase := NewGetSteamStoreDetailsUseCase(provider)

	result, err := useCase.Execute(context.Background(), "620")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	if result.Title != "Portal 2" {
		t.Fatalf("expected Portal 2, got %s", result.Title)
	}
}
