package integration

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"your-next-game-backend/internal/identity/domain"
	identitypg "your-next-game-backend/internal/identity/infrastructure/postgres"
	"your-next-game-backend/internal/platform/migrations"

	_ "github.com/lib/pq"
)

func TestIdentityPostgresRepository(t *testing.T) {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_DSN not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("ping db failed: %v", err)
	}

	migrationsDir := filepath.Clean("../../migrations")
	if err := migrations.ApplyUp(context.Background(), db, migrationsDir); err != nil {
		t.Fatalf("apply migrations failed: %v", err)
	}

	userID := mustUUID(t)
	// Generate random 17-digit SteamID
	steamID := fmt.Sprintf("7656%013d", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000000))

	repo := identitypg.NewUserRepository(db)
	ctx := context.Background()

	// Ensure user doesn't exist initially
	existing, err := repo.GetBySteamID(ctx, steamID)
	if err != nil {
		t.Fatalf("GetBySteamID failed: %v", err)
	}
	if existing != nil {
		t.Fatalf("expected user to not exist, got: %+v", existing)
	}

	newUser := domain.User{
		ID:         userID,
		SteamID:    steamID,
		Name:       "Test Steam Player",
		AvatarFull: "https://avatars.steamstatic.com/test_avatar.jpg",
	}

	// 1. Test Save (Insert)
	err = repo.Save(ctx, newUser)
	if err != nil {
		t.Fatalf("Save (insert) failed: %v", err)
	}

	// Defer cleanup of created rows
	defer func() {
		_, _ = db.Exec(`DELETE FROM steam_accounts WHERE user_id::text = $1`, userID)
		_, _ = db.Exec(`DELETE FROM users WHERE id::text = $1`, userID)
	}()

	// 2. Test GetBySteamID
	fetched, err := repo.GetBySteamID(ctx, steamID)
	if err != nil {
		t.Fatalf("GetBySteamID failed: %v", err)
	}
	if fetched == nil {
		t.Fatal("expected user to be found, got nil")
	}

	if fetched.ID != newUser.ID {
		t.Errorf("expected ID %q, got %q", newUser.ID, fetched.ID)
	}
	if fetched.SteamID != newUser.SteamID {
		t.Errorf("expected SteamID %q, got %q", newUser.SteamID, fetched.SteamID)
	}
	if fetched.Name != newUser.Name {
		t.Errorf("expected Name %q, got %q", newUser.Name, fetched.Name)
	}
	if fetched.AvatarFull != newUser.AvatarFull {
		t.Errorf("expected AvatarFull %q, got %q", newUser.AvatarFull, fetched.AvatarFull)
	}

	// 3. Test Save (Update / Conflict)
	updatedUser := domain.User{
		ID:         userID,
		SteamID:    steamID,
		Name:       "Updated Persona Name",
		AvatarFull: "https://avatars.steamstatic.com/updated_avatar.jpg",
	}

	err = repo.Save(ctx, updatedUser)
	if err != nil {
		t.Fatalf("Save (update) failed: %v", err)
	}

	fetchedUpdated, err := repo.GetBySteamID(ctx, steamID)
	if err != nil {
		t.Fatalf("GetBySteamID after update failed: %v", err)
	}
	if fetchedUpdated == nil {
		t.Fatal("expected updated user to be found, got nil")
	}

	if fetchedUpdated.ID != updatedUser.ID {
		t.Errorf("expected ID %q, got %q", updatedUser.ID, fetchedUpdated.ID)
	}
	if fetchedUpdated.Name != updatedUser.Name {
		t.Errorf("expected updated Name %q, got %q", updatedUser.Name, fetchedUpdated.Name)
	}
	if fetchedUpdated.AvatarFull != updatedUser.AvatarFull {
		t.Errorf("expected updated AvatarFull %q, got %q", updatedUser.AvatarFull, fetchedUpdated.AvatarFull)
	}
}
