package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	catalogapp "your-next-game-backend/internal/catalog/application"
	"your-next-game-backend/internal/platform/steamapi"
)

type fakeSteamStoreProvider struct {
	details *steamapi.StoreGameDetails
	err     error
}

func (f *fakeSteamStoreProvider) GetStoreDetails(_ context.Context, _ string) (*steamapi.StoreGameDetails, error) {
	return f.details, f.err
}

func TestGetSteamAppDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("bad input", func(t *testing.T) {
		handler := NewHandler(catalogapp.NewGetSteamStoreDetailsUseCase(&fakeSteamStoreProvider{}))
		req := httptest.NewRequest(http.MethodGet, "/v1/catalog/steam/apps/", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		handler.GetSteamAppDetails(c)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		handler := NewHandler(catalogapp.NewGetSteamStoreDetailsUseCase(&fakeSteamStoreProvider{err: context.DeadlineExceeded}))
		req := httptest.NewRequest(http.MethodGet, "/v1/catalog/steam/apps/620", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "appId", Value: "620"}}

		handler.GetSteamAppDetails(c)

		if rec.Code != http.StatusGatewayTimeout {
			t.Fatalf("expected 504, got %d", rec.Code)
		}
	})

	t.Run("dependency unavailable", func(t *testing.T) {
		handler := NewHandler(catalogapp.NewGetSteamStoreDetailsUseCase(&fakeSteamStoreProvider{err: errors.New("steam unavailable")}))
		req := httptest.NewRequest(http.MethodGet, "/v1/catalog/steam/apps/620", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "appId", Value: "620"}}

		handler.GetSteamAppDetails(c)

		if rec.Code != http.StatusBadGateway {
			t.Fatalf("expected 502, got %d", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		handler := NewHandler(catalogapp.NewGetSteamStoreDetailsUseCase(&fakeSteamStoreProvider{
			details: &steamapi.StoreGameDetails{
				Name: "Portal 2",
			},
		}))
		req := httptest.NewRequest(http.MethodGet, "/v1/catalog/steam/apps/620", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req
		c.Params = []gin.Param{{Key: "appId", Value: "620"}}

		handler.GetSteamAppDetails(c)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})
}
