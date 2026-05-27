package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	catalogapp "your-next-game-backend/internal/catalog/application"
	cataloghttp "your-next-game-backend/internal/catalog/interfaces/http"
	identityapp "your-next-game-backend/internal/identity/application"
	identitypg "your-next-game-backend/internal/identity/infrastructure/postgres"
	identitysteam "your-next-game-backend/internal/identity/infrastructure/steam"
	identityhttp "your-next-game-backend/internal/identity/interfaces/http"
	libraryapp "your-next-game-backend/internal/library/application"
	librarypg "your-next-game-backend/internal/library/infrastructure/postgres"
	libraryhttp "your-next-game-backend/internal/library/interfaces/http"
	"your-next-game-backend/internal/platform/config"
	"your-next-game-backend/internal/platform/httpserver"
	"your-next-game-backend/internal/platform/logger"
	"your-next-game-backend/internal/platform/steamapi"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	logger.InitLogger(cfg.AppVersion) // 'production' ou similar inicializará JSON


	db, dbReadiness, err := buildDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %v", err)
	}

	steamClient := steamapi.NewClient(cfg.SteamAPIKey, cfg.OutboundHTTPTimeout)
	jwtService := identityapp.NewJWTService(cfg.JWTSecret)

	// Identity Context
	userRepo := identitypg.NewUserRepository(db)
	createUserUseCase := identityapp.NewCreateUserUseCase(userRepo)
	steamLoginURLService := identityapp.NewSteamLoginURLService(identityapp.SteamLoginURLConfig{
		FrontendBaseURL: cfg.FrontendBaseURL,
		OpenIDRealm:     cfg.SteamOpenIDRealm,
		OpenIDReturnURL: cfg.SteamOpenIDReturnURL,
	})
	openidValidator := identitysteam.NewOpenIDValidator()
	steamExchangeService := identityapp.NewSteamExchangeService(steamClient, openidValidator)
	identityHandler := identityhttp.NewHandler(
		steamLoginURLService,
		steamExchangeService,
		jwtService,
		createUserUseCase,
		userRepo,
		cfg.FrontendSuccessURL,
	)

	// Library Context
	gamesRepo := librarypg.NewRepository(db)
	listGamesUseCase := libraryapp.NewListGamesUseCase(gamesRepo)
	getGameByIDUseCase := libraryapp.NewGetGameByIDUseCase(gamesRepo)
	listSteamGamesUseCase := libraryapp.NewListSteamGamesUseCase(steamClient)
	createGameUseCase := libraryapp.NewCreateGameUseCase(gamesRepo)
	libraryHandler := libraryhttp.NewHandler(
		listGamesUseCase,
		getGameByIDUseCase,
		listSteamGamesUseCase,
		createGameUseCase,
	)

	// Catalog Context
	catalogSteamStoreDetailsUseCase := catalogapp.NewGetSteamStoreDetailsUseCase(steamClient)
	catalogHandler := cataloghttp.NewHandler(catalogSteamStoreDetailsUseCase)

	server := &http.Server{
		Addr: ":" + cfg.HTTPPort,
		Handler: httpserver.New(httpserver.Dependencies{
			Identity:  identityHandler,
			Library:   libraryHandler,
			Catalog:   catalogHandler,
			TokenSvc:  jwtService,
			Readiness: dbReadiness,
			Version:   cfg.AppVersion,
		}),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	log.Printf("your-next-game-backend listening on :%s", cfg.HTTPPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func buildDB(cfg config.Config) (*sql.DB, func(context.Context) error, error) {
	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres open failed: %w", err)
	}
	
	// Configure connection pooling
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("postgres ping failed: %w", err)
	}
	log.Printf("connected to postgres")

	readiness := func(ctx context.Context) error {
		return db.PingContext(ctx)
	}

	return db, readiness, nil
}
