package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	libraryapp "your-next-game-backend/internal/library/application"
	"your-next-game-backend/internal/library/domain"
	"your-next-game-backend/internal/platform/steamapi"
)

type fakeGameRepo struct {
	games []domain.Game
	err   error
}

func (f *fakeGameRepo) ListByUserID(_ context.Context, _ string) ([]domain.Game, error) {
	return f.games, f.err
}

func (f *fakeGameRepo) GetByUserIDAndGameID(_ context.Context, _, _ string) (*domain.Game, error) {
	if f.err != nil {
		return nil, f.err
	}
	if len(f.games) > 0 {
		return &f.games[0], nil
	}
	return nil, nil
}

func (f *fakeGameRepo) CreateGame(_ context.Context, _ string, game domain.Game) error {
	f.games = append(f.games, game)
	return f.err
}

type fakeSteamProfileProvider struct {
	err error
}

func (f *fakeSteamProfileProvider) GetOwnedGames(_ context.Context, _ string) ([]steamapi.OwnedGame, error) {
	return nil, f.err
}

func TestListGamesByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing user id", func(t *testing.T) {
		repo := &fakeGameRepo{}
		handler := NewHandler(
			libraryapp.NewListGamesUseCase(repo),
			libraryapp.NewGetGameByIDUseCase(repo),
			libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{}),
			libraryapp.NewCreateGameUseCase(repo),
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/users//games", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "userId", Value: ""}}

		handler.ListGamesByUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		repo := &fakeGameRepo{err: context.DeadlineExceeded}
		handler := NewHandler(
			libraryapp.NewListGamesUseCase(repo),
			libraryapp.NewGetGameByIDUseCase(repo),
			libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{}),
			libraryapp.NewCreateGameUseCase(repo),
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/games", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "userId", Value: "user-1"}}

		handler.ListGamesByUser(c)

		if rec.Code != http.StatusGatewayTimeout {
			t.Fatalf("expected 504, got %d", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		repo := &fakeGameRepo{
			games: []domain.Game{{ID: "1", Title: "Game 1"}},
		}
		handler := NewHandler(
			libraryapp.NewListGamesUseCase(repo),
			libraryapp.NewGetGameByIDUseCase(repo),
			libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{}),
			libraryapp.NewCreateGameUseCase(repo),
		)

		req := httptest.NewRequest(http.MethodGet, "/v1/users/user-1/games", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "userId", Value: "user-1"}}

		handler.ListGamesByUser(c)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var payload []domain.Game
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode failed: %v", err)
		}

		if len(payload) != 1 {
			t.Fatal("expected 1 game")
		}
	})
}

func TestGetGameByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeGameRepo{
		games: []domain.Game{{ID: "g1", Title: "Game"}},
	}
	handler := NewHandler(
		libraryapp.NewListGamesUseCase(repo),
		libraryapp.NewGetGameByIDUseCase(repo),
		libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{}),
		libraryapp.NewCreateGameUseCase(repo),
	)

	t.Run("bad input", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/users/u1/games/", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "userId", Value: "u1"}, {Key: "gameId", Value: ""}}
		
		handler.GetGameByUser(c)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		timeoutHandler := NewHandler(
			libraryapp.NewListGamesUseCase(repo),
			libraryapp.NewGetGameByIDUseCase(&fakeGameRepo{err: context.DeadlineExceeded}),
			libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{}),
			libraryapp.NewCreateGameUseCase(repo),
		)
		req := httptest.NewRequest(http.MethodGet, "/v1/users/u1/games/g1", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "userId", Value: "u1"}, {Key: "gameId", Value: "g1"}}
		
		timeoutHandler.GetGameByUser(c)
		if rec.Code != http.StatusGatewayTimeout {
			t.Fatalf("expected 504, got %d", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/users/u1/games/g1", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "userId", Value: "u1"}, {Key: "gameId", Value: "g1"}}
		
		handler.GetGameByUser(c)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})
}

func TestListSteamGames(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(
		libraryapp.NewListGamesUseCase(&fakeGameRepo{}),
		libraryapp.NewGetGameByIDUseCase(&fakeGameRepo{}),
		libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{}),
		libraryapp.NewCreateGameUseCase(&fakeGameRepo{}),
	)

	t.Run("bad input", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/integrations/steam/users//games", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "steamId", Value: ""}}

		handler.ListSteamGames(c)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		timeoutHandler := NewHandler(
			libraryapp.NewListGamesUseCase(&fakeGameRepo{}),
			libraryapp.NewGetGameByIDUseCase(&fakeGameRepo{}),
			libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{err: context.DeadlineExceeded}),
			libraryapp.NewCreateGameUseCase(&fakeGameRepo{}),
		)
		req := httptest.NewRequest(http.MethodGet, "/v1/integrations/steam/users/s1/games", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "steamId", Value: "s1"}}

		timeoutHandler.ListSteamGames(c)
		if rec.Code != http.StatusGatewayTimeout {
			t.Fatalf("expected 504, got %d", rec.Code)
		}
	})

	t.Run("dependency unavailable", func(t *testing.T) {
		failHandler := NewHandler(
			libraryapp.NewListGamesUseCase(&fakeGameRepo{}),
			libraryapp.NewGetGameByIDUseCase(&fakeGameRepo{}),
			libraryapp.NewListSteamGamesUseCase(&fakeSteamProfileProvider{err: errors.New("steam unavailable")}),
			libraryapp.NewCreateGameUseCase(&fakeGameRepo{}),
		)
		req := httptest.NewRequest(http.MethodGet, "/v1/integrations/steam/users/s1/games", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "steamId", Value: "s1"}}

		failHandler.ListSteamGames(c)
		if rec.Code != http.StatusBadGateway {
			t.Fatalf("expected 502, got %d", rec.Code)
		}
	})
}
