package application

import (
	"context"
	"math"

	"your-next-game-backend/internal/library/domain"
	"your-next-game-backend/internal/platform/steamapi"
)

type SteamOwnedGamesProvider interface {
	GetOwnedGames(ctx context.Context, steamID string) ([]steamapi.OwnedGame, error)
}

type ListSteamGamesUseCase struct {
	provider SteamOwnedGamesProvider
}

func NewListSteamGamesUseCase(provider SteamOwnedGamesProvider) *ListSteamGamesUseCase {
	return &ListSteamGamesUseCase{
		provider: provider,
	}
}

func (u *ListSteamGamesUseCase) Execute(ctx context.Context, steamID string) ([]domain.Game, error) {
	ownedGames, err := u.provider.GetOwnedGames(ctx, steamID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Game, 0, len(ownedGames))
	for _, game := range ownedGames {
		appID := steamapi.AppIDToString(game.AppID)
		hours := math.Round((float64(game.PlaytimeForever)/60)*10) / 10

		status := "Não Jogado"
		if game.PlaytimeForever > 1 {
			status = "Jogado"
		}

		result = append(result, domain.Game{
			ID:          appID,
			SteamAppID:  &appID,
			Title:       game.Name,
			CoverURL:    "https://cdn.akamai.steamstatic.com/steam/apps/" + appID + "/library_600x900_2x.jpg",
			Status:      status,
			Genre:       "Steam",
			Platform:    "Steam",
			HoursPlayed: hours,
			Description: "",
			Developer:   "",
			ReleaseDate: "",
			Rating:      nil,
		})
	}

	return result, nil
}
