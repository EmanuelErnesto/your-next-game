package application

import (
	"context"

	"your-next-game-backend/internal/platform/steamapi"
)

type SteamProfileProvider interface {
	GetPlayerSummary(ctx context.Context, steamID string) (*steamapi.PlayerSummary, error)
}

type OpenIDValidator interface {
	Validate(ctx context.Context, params map[string]string) (string, error)
}

type SteamExchangeService struct {
	profileProvider SteamProfileProvider
	openidValidator OpenIDValidator
}

type AuthUser struct {
	ID      string  `json:"id"`
	SteamID string  `json:"steamId"`
	Name    string  `json:"name"`
	Image   string  `json:"image"`
	Email   *string `json:"email"`
}

func NewSteamExchangeService(profileProvider SteamProfileProvider, openidValidator OpenIDValidator) *SteamExchangeService {
	return &SteamExchangeService{
		profileProvider: profileProvider,
		openidValidator: openidValidator,
	}
}

func (s *SteamExchangeService) Exchange(ctx context.Context, params map[string]string) (*AuthUser, error) {
	steamID, err := s.openidValidator.Validate(ctx, params)
	if err != nil {
		return nil, err
	}
	if steamID == "" {
		return nil, nil
	}

	profile, err := s.profileProvider.GetPlayerSummary(ctx, steamID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, nil
	}

	return &AuthUser{
		ID:      steamID,
		SteamID: steamID,
		Name:    profile.PersonaName,
		Image:   profile.AvatarFull,
		Email:   nil,
	}, nil
}
