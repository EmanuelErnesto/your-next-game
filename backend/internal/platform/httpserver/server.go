package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	cataloghttp "your-next-game-backend/internal/catalog/interfaces/http"
	identityapp "your-next-game-backend/internal/identity/application"
	identityhttp "your-next-game-backend/internal/identity/interfaces/http"
	libraryhttp "your-next-game-backend/internal/library/interfaces/http"
)

type Dependencies struct {
	Identity   *identityhttp.Handler
	Library    *libraryhttp.Handler
	Catalog    *cataloghttp.Handler
	TokenSvc   identityapp.TokenService
	Readiness  func(ctx context.Context) error
	Version    string
}

func New(dep Dependencies) http.Handler {
	router := gin.New() // usando gin.New() em vez de Default para evitar o logger padrão não-estruturado do Gin.
	router.Use(LoggerMiddleware(), gin.Recovery())

	// Health and Readiness
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		if dep.Readiness != nil {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
			defer cancel()

			if err := dep.Readiness(ctx); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "not_ready",
					"error":  "dependency_unavailable",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	api := router.Group("/v1")
	{
		api.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": dep.Version})
		})

		// Identity / Auth Routes
		authGroup := api.Group("/auth")
		{
			authGroup.GET("/steam/login", dep.Identity.Login)
			authGroup.GET("/steam/callback", dep.Identity.SteamCallback)
			authGroup.GET("/logout", dep.Identity.Logout)
		}

		// Public catalog routes
		api.GET("/catalog/steam/apps/:appId", dep.Catalog.GetSteamAppDetails)

		// Protected Routes
		protected := api.Group("")
		protected.Use(authMiddleware(dep.TokenSvc))
		{
			protected.GET("/auth/me", dep.Identity.Me)
			users := protected.Group("/users/:userId")
			{
				users.GET("/games", dep.Library.ListGamesByUser)
				users.POST("/games", dep.Library.PostCreateGame)
				users.GET("/games/:gameId", dep.Library.GetGameByUser)
			}
			protected.GET("/integrations/steam/users/:steamId/games", dep.Library.ListSteamGames)
		}
	}

	return router
}

func authMiddleware(tokenSvc identityapp.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := tokenSvc.VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
