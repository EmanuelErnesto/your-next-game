package application

import (
	"context"
	"strings"

	"your-next-game-backend/internal/platform/steamapi"
)

type SteamStoreDetailsProvider interface {
	GetStoreDetails(ctx context.Context, appID string) (*steamapi.StoreGameDetails, error)
}

type StoreGameProjection struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Developer   string `json:"developer"`
	Genre       string `json:"genre"`
	ReleaseDate string `json:"releaseDate"`
	CoverURL    string `json:"coverUrl"`
}

type GetSteamStoreDetailsUseCase struct {
	provider SteamStoreDetailsProvider
}

func NewGetSteamStoreDetailsUseCase(provider SteamStoreDetailsProvider) *GetSteamStoreDetailsUseCase {
	return &GetSteamStoreDetailsUseCase{
		provider: provider,
	}
}

func (u *GetSteamStoreDetailsUseCase) Execute(
	ctx context.Context,
	appID string,
) (*StoreGameProjection, error) {
	details, err := u.provider.GetStoreDetails(ctx, appID)
	if err != nil {
		return nil, err
	}
	if details == nil {
		return nil, nil
	}

	genre := "Steam"
	if len(details.Genres) > 0 {
		genres := make([]string, 0, len(details.Genres))
		for _, g := range details.Genres {
			if g.Description != "" {
				genres = append(genres, g.Description)
			}
		}
		if len(genres) > 0 {
			genre = strings.Join(genres, ", ")
		}
	}

	developer := ""
	if len(details.Developers) > 0 {
		developer = details.Developers[0]
	}

	description := details.DetailedDescription
	if description == "" {
		description = details.ShortDescription
	}

	cover := details.HeaderImage
	if cover == "" {
		cover = details.Background
	}

	return &StoreGameProjection{
		Title:       details.Name,
		Description: description,
		Developer:   developer,
		Genre:       genre,
		ReleaseDate: details.ReleaseDate.Date,
		CoverURL:    cover,
	}, nil
}
