package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	identityapp "your-next-game-backend/internal/identity/application"
	identitydomain "your-next-game-backend/internal/identity/domain"
	"your-next-game-backend/internal/platform/logger"
)

type Handler struct {
	loginURL      *identityapp.SteamLoginURLService
	exchangeSteam *identityapp.SteamExchangeService
	jwtService    identityapp.TokenService
	createUser    *identityapp.CreateUserUseCase
	userRepo      identitydomain.UserRepository
	frontendURL   string
}

func NewHandler(
	loginURL *identityapp.SteamLoginURLService,
	exchangeSteam *identityapp.SteamExchangeService,
	jwtService identityapp.TokenService,
	createUser *identityapp.CreateUserUseCase,
	userRepo identitydomain.UserRepository,
	frontendURL string,
) *Handler {
	return &Handler{
		loginURL:      loginURL,
		exchangeSteam: exchangeSteam,
		jwtService:    jwtService,
		createUser:    createUser,
		userRepo:      userRepo,
		frontendURL:   frontendURL,
	}
}

func (h *Handler) Login(c *gin.Context) {
	c.Redirect(http.StatusFound, h.loginURL.Build())
}

func (h *Handler) SteamCallback(c *gin.Context) {
	// Extract all query parameters into a map
	queryParams := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			queryParams[k] = v[0]
		}
	}

	user, err := h.exchangeSteam.Exchange(c.Request.Context(), queryParams)
	if err != nil {
		logger.ContextLogger(c.Request.Context()).Error("Steam exchange failed", "error", err)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "steam exchange timeout"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "steam exchange failed"})
		return
	}
	
	if user == nil {
		logger.ContextLogger(c.Request.Context()).Warn("Steam exchange returned nil user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create user in DB
	createdUser, err := h.createUser.Execute(c.Request.Context(), identitydomain.User{
		ID:         user.ID,
		SteamID:    user.ID,
		Name:       user.Name,
		AvatarFull: user.Image,
	})
	if err != nil {
		logger.ContextLogger(c.Request.Context()).Error("Create user in DB failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	accessToken, refreshToken, err := h.jwtService.GenerateTokens(createdUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	secure := strings.HasPrefix(h.frontendURL, "https://")
	c.SetCookie("access_token", accessToken, 15*60, "/", "", secure, true)
	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", secure, true)

	c.Redirect(http.StatusFound, h.frontendURL)
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if h.userRepo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user repository unavailable"})
		return
	}

	user, err := h.userRepo.GetBySteamID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, identityapp.AuthUser{
		ID:      user.ID,
		SteamID: user.SteamID,
		Name:    user.Name,
		Image:   user.AvatarFull,
		Email:   nil,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	secure := strings.HasPrefix(h.frontendURL, "https://")
	c.SetCookie("access_token", "", -1, "/", "", secure, true)
	c.SetCookie("refresh_token", "", -1, "/", "", secure, true)
	c.Redirect(http.StatusFound, h.frontendURL)
}
