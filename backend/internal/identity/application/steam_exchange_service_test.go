package application

import (
	"context"
	"errors"
	"testing"

	"your-next-game-backend/internal/platform/steamapi"
)

type fakeSteamProfileProvider struct {
	err error
}

func (f *fakeSteamProfileProvider) GetPlayerSummary(_ context.Context, _ string) (*steamapi.PlayerSummary, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &steamapi.PlayerSummary{PersonaName: "Player", AvatarFull: "image"}, nil
}

type fakeOpenIDValidator struct {
	steamID string
	err     error
}

func (f *fakeOpenIDValidator) Validate(ctx context.Context, params map[string]string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.steamID, nil
}

func TestSteamExchangeService(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := NewSteamExchangeService(&fakeSteamProfileProvider{}, &fakeOpenIDValidator{steamID: "123"})
		user, err := service.Exchange(context.Background(), map[string]string{"openid.claimed_id": "123"})
		if err != nil {
			t.Fatalf("expected nil err, got %v", err)
		}
		if user == nil || user.ID != "123" {
			t.Fatal("expected user with ID 123")
		}
	})

	t.Run("validator error", func(t *testing.T) {
		service := NewSteamExchangeService(&fakeSteamProfileProvider{}, &fakeOpenIDValidator{err: errors.New("validator err")})
		_, err := service.Exchange(context.Background(), map[string]string{})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("profile provider error", func(t *testing.T) {
		service := NewSteamExchangeService(&fakeSteamProfileProvider{err: errors.New("provider err")}, &fakeOpenIDValidator{steamID: "123"})
		_, err := service.Exchange(context.Background(), map[string]string{"openid.claimed_id": "123"})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
