package integration

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	librarypg "your-next-game-backend/internal/library/infrastructure/postgres"
	"your-next-game-backend/internal/platform/migrations"

	_ "github.com/lib/pq"
)

func TestPostgresRepository_ListAndGetByUserID(t *testing.T) {
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
	catalogID := mustUUID(t)
	userGameID := mustUUID(t)

	_, err = db.Exec(`
insert into users (id, display_name, avatar_url)
values ($1, 'Integration User', '')
`, userID)
	if err != nil {
		t.Fatalf("insert users failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(`delete from users where id::text = $1`, userID)
	}()

	_, err = db.Exec(`
insert into games_catalog (id, provider, provider_game_id, title, cover_url, description, developer, genre, release_date)
values ($1, 'steam', '999999', 'Integration Game', 'https://example.com/cover.jpg', 'Description', 'Dev', 'RPG', '2020-01-01')
`, catalogID)
	if err != nil {
		t.Fatalf("insert games_catalog failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(`delete from games_catalog where id::text = $1`, catalogID)
	}()

	_, err = db.Exec(`
insert into user_games (id, user_id, catalog_game_id, steam_app_id, status, hours_played, rating, source)
values ($1, $2, $3, 999999, 'not_played', 12.5, 4.5, 'manual')
`, userGameID, userID, catalogID)
	if err != nil {
		t.Fatalf("insert user_games failed: %v", err)
	}
	defer func() {
		_, _ = db.Exec(`delete from user_games where id::text = $1`, userGameID)
	}()

	repo := librarypg.NewRepository(db)

	games, err := repo.ListByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("list by user id failed: %v", err)
	}
	if len(games) == 0 {
		t.Fatal("expected at least one game")
	}

	game, err := repo.GetByUserIDAndGameID(context.Background(), userID, userGameID)
	if err != nil {
		t.Fatalf("get by user/game id failed: %v", err)
	}
	if game == nil {
		t.Fatal("expected game, got nil")
	}
	if game.ID != userGameID {
		t.Fatalf("expected game id %s, got %s", userGameID, game.ID)
	}
	if game.Title != "Integration Game" {
		t.Fatalf("unexpected title: %s", game.Title)
	}
}

func mustUUID(t *testing.T) string {
	t.Helper()

	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		t.Fatalf("uuid random read failed: %v", err)
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	)
}
