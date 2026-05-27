package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	catalogapp "your-next-game-backend/internal/catalog/application"
	cataloghttp "your-next-game-backend/internal/catalog/interfaces/http"
	identityapp "your-next-game-backend/internal/identity/application"
	identitydomain "your-next-game-backend/internal/identity/domain"
	identityhttp "your-next-game-backend/internal/identity/interfaces/http"
	libraryapp "your-next-game-backend/internal/library/application"
	librarydomain "your-next-game-backend/internal/library/domain"
	libraryhttp "your-next-game-backend/internal/library/interfaces/http"
	"your-next-game-backend/internal/platform/steamapi"
)

type stubRepository struct {
	games []librarydomain.Game
}

type stubUserRepository struct {
	user *identitydomain.User
}

func (r *stubRepository) ListByUserID(context.Context, string) ([]librarydomain.Game, error) {
	result := make([]librarydomain.Game, len(r.games))
	copy(result, r.games)
	return result, nil
}

func (r *stubRepository) GetByUserIDAndGameID(
	_ context.Context,
	_ string,
	gameID string,
) (*librarydomain.Game, error) {
	for _, game := range r.games {
		if game.ID == gameID {
			g := game
			return &g, nil
		}
	}
	return nil, nil
}

func (r *stubRepository) CreateGame(_ context.Context, _ string, game librarydomain.Game) error {
	r.games = append(r.games, game)
	return nil
}

func (r *stubUserRepository) Save(_ context.Context, user identitydomain.User) error {
	r.user = &user
	return nil
}

func (r *stubUserRepository) GetBySteamID(_ context.Context, steamID string) (*identitydomain.User, error) {
	if r.user != nil && r.user.SteamID == steamID {
		return r.user, nil
	}
	return nil, nil
}

func TestServerRoutes(t *testing.T) {
	steamClient := steamapi.NewClient("", 0)

	tokenSvc := identityapp.NewJWTService("test_secret")

	identityHandler := identityhttp.NewHandler(
		identityapp.NewSteamLoginURLService(identityapp.SteamLoginURLConfig{
			FrontendBaseURL: "http://localhost:3001",
			OpenIDRealm:     "http://localhost:3001",
			OpenIDReturnURL: "http://localhost:3001/api/auth/steam/callback",
		}),
		identityapp.NewSteamExchangeService(steamClient, nil),
		tokenSvc,
		nil,
		&stubUserRepository{},
		"http://localhost:3001",
	)

	repo := &stubRepository{
		games: []librarydomain.Game{
			{
				ID:       "local-1",
				Title:    "Demo Game",
				CoverURL: "https://example.com/cover.jpg",
				Status:   "Não Jogado",
				Genre:    "RPG",
				Platform: "Steam",
			},
		},
	}
	libraryHandler := libraryhttp.NewHandler(
		libraryapp.NewListGamesUseCase(repo),
		libraryapp.NewGetGameByIDUseCase(repo),
		libraryapp.NewListSteamGamesUseCase(steamClient),
		libraryapp.NewCreateGameUseCase(repo),
	)

	catalogHandler := cataloghttp.NewHandler(catalogapp.NewGetSteamStoreDetailsUseCase(steamClient))

	handler := New(Dependencies{
		Identity:  identityHandler,
		Library:   libraryHandler,
		Catalog:   catalogHandler,
		TokenSvc:  tokenSvc,
		Readiness: func(context.Context) error { return nil },
		Version:   "test",
	})

	t.Run("healthz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("list games", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/games", nil)
		// Set cookie
		acc, _, _ := tokenSvc.GenerateTokens("user-1")
		req.AddCookie(&http.Cookie{
			Name:     "access_token",
			Value:    acc,
			Expires:  time.Now().Add(15 * time.Minute),
			HttpOnly: true,
			Path:     "/",
		})

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var payload []map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if len(payload) == 0 {
			t.Fatal("expected at least one game")
		}
	})

	t.Run("readyz unhealthy", func(t *testing.T) {
		unhealthyHandler := New(Dependencies{
			Identity: identityHandler,
			Library:  libraryHandler,
			Catalog:  catalogHandler,
			TokenSvc: tokenSvc,
			Readiness: func(context.Context) error {
				return errors.New("db unavailable")
			},
			Version: "test",
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()

		unhealthyHandler.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d", rec.Code)
		}
	})
}
