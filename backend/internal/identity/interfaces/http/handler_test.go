package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"your-next-game-backend/internal/identity/application"
	"your-next-game-backend/internal/identity/domain"
	"your-next-game-backend/internal/platform/steamapi"
)

type fakeSteamProfileProvider struct {
	err error
}

func (f *fakeSteamProfileProvider) GetPlayerSummary(_ context.Context, _ string) (*steamapi.PlayerSummary, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &steamapi.PlayerSummary{
		PersonaName: "Player",
		AvatarFull:  "https://image",
	}, nil
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

type fakeTokenSvc struct {
	access  string
	refresh string
	err     error
}

func (f *fakeTokenSvc) GenerateTokens(_ string) (string, string, error) {
	if f.err != nil {
		return "", "", f.err
	}
	return f.access, f.refresh, nil
}

func (f *fakeTokenSvc) VerifyToken(tokenString string) (string, error) {
	return "user-1", nil
}

type fakeUserRepository struct {
	user domain.User
	err  error
}

func (f *fakeUserRepository) GetBySteamID(ctx context.Context, steamID string) (*domain.User, error) {
	if f.user.SteamID == steamID {
		return &f.user, nil
	}
	return nil, nil
}

func (f *fakeUserRepository) Save(ctx context.Context, user domain.User) error {
	f.user = user
	return f.err
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	loginURLSvc := application.NewSteamLoginURLService(application.SteamLoginURLConfig{
		FrontendBaseURL: "http://frontend",
		OpenIDRealm:     "http://backend",
		OpenIDReturnURL: "http://backend/callback",
	})
	handler := NewHandler(loginURLSvc, nil, nil, nil, nil, "http://frontend")

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/steam/login", nil)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	handler.Login(c)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc == "" {
		t.Fatal("expected location header")
	}
}

func TestSteamCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	loginURLSvc := application.NewSteamLoginURLService(application.SteamLoginURLConfig{})

	t.Run("success http secure false", func(t *testing.T) {
		exchangeSvc := application.NewSteamExchangeService(&fakeSteamProfileProvider{}, &fakeOpenIDValidator{steamID: "123"})
		tokenSvc := &fakeTokenSvc{access: "acc", refresh: "ref"}
		createUser := application.NewCreateUserUseCase(&fakeUserRepository{})
		handler := NewHandler(loginURLSvc, exchangeSvc, tokenSvc, createUser, &fakeUserRepository{}, "http://frontend")

		req := httptest.NewRequest(http.MethodGet, "/v1/auth/steam/callback?openid.mode=id_res", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		handler.SteamCallback(c)

		if rec.Code != http.StatusFound {
			t.Fatalf("expected 302, got %d", rec.Code)
		}

		cookies := rec.Result().Cookies()
		if len(cookies) != 2 {
			t.Fatalf("expected 2 cookies, got %d", len(cookies))
		}
		for _, cookie := range cookies {
			if cookie.Secure {
				t.Errorf("expected cookie %s Secure to be false on http, got true", cookie.Name)
			}
		}
	})

	t.Run("success https secure true", func(t *testing.T) {
		exchangeSvc := application.NewSteamExchangeService(&fakeSteamProfileProvider{}, &fakeOpenIDValidator{steamID: "123"})
		tokenSvc := &fakeTokenSvc{access: "acc", refresh: "ref"}
		createUser := application.NewCreateUserUseCase(&fakeUserRepository{})
		handler := NewHandler(loginURLSvc, exchangeSvc, tokenSvc, createUser, &fakeUserRepository{}, "https://frontend")

		req := httptest.NewRequest(http.MethodGet, "/v1/auth/steam/callback?openid.mode=id_res", nil)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		handler.SteamCallback(c)

		if rec.Code != http.StatusFound {
			t.Fatalf("expected 302, got %d", rec.Code)
		}

		cookies := rec.Result().Cookies()
		if len(cookies) != 2 {
			t.Fatalf("expected 2 cookies, got %d", len(cookies))
		}
		for _, cookie := range cookies {
			if !cookie.Secure {
				t.Errorf("expected cookie %s Secure to be true on https, got false", cookie.Name)
			}
		}
	})
}

func TestMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(nil, nil, nil, nil, &fakeUserRepository{
		user: domain.User{
			ID:         "123",
			SteamID:    "123",
			Name:       "Player",
			AvatarFull: "https://image",
		},
	}, "http://frontend")

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/me", nil)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Set("userID", "123")

	handler.Me(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(nil, nil, nil, nil, nil, "http://frontend")

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/logout", nil)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	handler.Logout(c)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
}
