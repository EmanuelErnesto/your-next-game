package http

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	libraryapp "your-next-game-backend/internal/library/application"
	librarydomain "your-next-game-backend/internal/library/domain"
)

type Handler struct {
	listGames      *libraryapp.ListGamesUseCase
	getGameByID    *libraryapp.GetGameByIDUseCase
	listSteamGames *libraryapp.ListSteamGamesUseCase
	createGame     *libraryapp.CreateGameUseCase
}

func NewHandler(
	listGames *libraryapp.ListGamesUseCase,
	getGameByID *libraryapp.GetGameByIDUseCase,
	listSteamGames *libraryapp.ListSteamGamesUseCase,
	createGame *libraryapp.CreateGameUseCase,
) *Handler {
	return &Handler{
		listGames:      listGames,
		getGameByID:    getGameByID,
		listSteamGames: listSteamGames,
		createGame:     createGame,
	}
}

func (h *Handler) ListGamesByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	games, err := h.listGames.Execute(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "list games timeout"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (h *Handler) GetGameByUser(c *gin.Context) {
	userID := c.Param("userId")
	gameID := c.Param("gameId")
	if userID == "" || gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and gameId are required"})
		return
	}

	game, err := h.getGameByID.Execute(c.Request.Context(), userID, gameID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "get game timeout"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get game"})
		return
	}

	c.JSON(http.StatusOK, game)
}

func (h *Handler) ListSteamGames(c *gin.Context) {
	steamID := c.Param("steamId")
	if steamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "steamId is required"})
		return
	}

	games, err := h.listSteamGames.Execute(c.Request.Context(), steamID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "list steam games timeout"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to list steam games"})
		return
	}

	c.JSON(http.StatusOK, games)
}

func (h *Handler) PostCreateGame(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	var game librarydomain.Game
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if game.ID == "" || game.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID and title are required"})
		return
	}

	err := h.createGame.Execute(c.Request.Context(), userID, game)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "create game timeout"})
			return
		}
		log.Printf("CREATE GAME ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}
