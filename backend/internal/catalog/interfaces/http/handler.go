package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	catalogapp "your-next-game-backend/internal/catalog/application"
)

type Handler struct {
	getSteamStoreDetails *catalogapp.GetSteamStoreDetailsUseCase
}

func NewHandler(getSteamStoreDetails *catalogapp.GetSteamStoreDetailsUseCase) *Handler {
	return &Handler{
		getSteamStoreDetails: getSteamStoreDetails,
	}
}

func (h *Handler) GetSteamAppDetails(c *gin.Context) {
	appID := c.Param("appId")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "appId is required"})
		return
	}

	details, err := h.getSteamStoreDetails.Execute(c.Request.Context(), appID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "steam app details timeout"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to get steam app details"})
		return
	}

	c.JSON(http.StatusOK, details)
}
