package application

import (
	"strings"
	"testing"
)

func TestSteamLoginURLServiceBuild(t *testing.T) {
	service := NewSteamLoginURLService(SteamLoginURLConfig{
		FrontendBaseURL: "http://localhost:3001",
		OpenIDRealm:     "http://localhost:3001",
		OpenIDReturnURL: "http://localhost:3001/api/auth/steam/callback",
	})

	url := service.Build()

	if !strings.HasPrefix(url, "https://steamcommunity.com/openid/login?") {
		t.Fatalf("expected steam login base URL, got %s", url)
	}
	if !strings.Contains(url, "openid.return_to=http%3A%2F%2Flocalhost%3A3001%2Fapi%2Fauth%2Fsteam%2Fcallback") {
		t.Fatalf("expected return_to param, got %s", url)
	}
}
