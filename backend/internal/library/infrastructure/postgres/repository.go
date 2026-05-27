package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"your-next-game-backend/internal/library/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListByUserID(ctx context.Context, userID string) ([]domain.Game, error) {
	const query = `
select
  ug.id,
  gc.title,
  coalesce(gc.cover_url, '') as cover_url,
  case ug.status
    when 'played' then 'Jogado'
    else 'Não Jogado'
  end as status,
  coalesce(gc.genre, 'Steam') as genre,
  coalesce(gc.provider, 'Steam') as platform,
  coalesce(gc.description, '') as description,
  coalesce(gc.developer, '') as developer,
  coalesce(to_char(gc.release_date, 'YYYY-MM-DD'), '') as release_date,
  ug.hours_played,
  ug.rating,
  case when ug.steam_app_id is null then null else ug.steam_app_id::text end as steam_app_id
from user_games ug
join games_catalog gc on gc.id = ug.catalog_game_id
where ug.user_id::text = $1
order by gc.title asc
`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list games query failed: %w", err)
	}
	defer rows.Close()

	games := make([]domain.Game, 0)
	for rows.Next() {
		game, err := scanGame(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list games rows failed: %w", err)
	}

	return games, nil
}

func (r *Repository) GetByUserIDAndGameID(
	ctx context.Context,
	userID, gameID string,
) (*domain.Game, error) {
	const query = `
select
  ug.id,
  gc.title,
  coalesce(gc.cover_url, '') as cover_url,
  case ug.status
    when 'played' then 'Jogado'
    else 'Não Jogado'
  end as status,
  coalesce(gc.genre, 'Steam') as genre,
  coalesce(gc.provider, 'Steam') as platform,
  coalesce(gc.description, '') as description,
  coalesce(gc.developer, '') as developer,
  coalesce(to_char(gc.release_date, 'YYYY-MM-DD'), '') as release_date,
  ug.hours_played,
  ug.rating,
  case when ug.steam_app_id is null then null else ug.steam_app_id::text end as steam_app_id
from user_games ug
join games_catalog gc on gc.id = ug.catalog_game_id
where ug.user_id::text = $1 and ug.id::text = $2
limit 1
`

	row := r.db.QueryRowContext(ctx, query, userID, gameID)
	game, err := scanGame(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &game, nil
}

type gameScanner interface {
	Scan(dest ...any) error
}

func scanGame(scanner gameScanner) (domain.Game, error) {
	var game domain.Game
	var rating sql.NullFloat64
	var steamAppID sql.NullString

	err := scanner.Scan(
		&game.ID,
		&game.Title,
		&game.CoverURL,
		&game.Status,
		&game.Genre,
		&game.Platform,
		&game.Description,
		&game.Developer,
		&game.ReleaseDate,
		&game.HoursPlayed,
		&rating,
		&steamAppID,
	)
	if err != nil {
		return domain.Game{}, fmt.Errorf("scan game failed: %w", err)
	}

	if rating.Valid {
		value := rating.Float64
		game.Rating = &value
	}
	if steamAppID.Valid {
		value := steamAppID.String
		game.SteamAppID = &value
	}

	return game, nil
}

func (r *Repository) CreateGame(ctx context.Context, userID string, game domain.Game) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	providerGameID := game.ID
	if game.SteamAppID != nil && *game.SteamAppID != "" {
		providerGameID = *game.SteamAppID
	}

	// Insert into games_catalog
	catalogQuery := `
		INSERT INTO games_catalog (id, provider, provider_game_id, title, cover_url, genre, description, developer, release_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (provider, provider_game_id) DO UPDATE SET title = EXCLUDED.title
		RETURNING id
	`
	var catalogGameID string
	err = tx.QueryRowContext(ctx, catalogQuery,
		game.ID,
		game.Platform,
		providerGameID,
		game.Title,
		game.CoverURL,
		game.Genre,
		game.Description,
		game.Developer,
		sql.NullString{String: game.ReleaseDate, Valid: game.ReleaseDate != ""},
	).Scan(&catalogGameID)
	if err != nil {
		return fmt.Errorf("failed to insert into games_catalog: %w", err)
	}

	// Insert into user_games
	userGamesQuery := `
		INSERT INTO user_games (id, user_id, catalog_game_id, status, hours_played, rating, steam_app_id, source)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, catalog_game_id) DO NOTHING
	`
	
	status := "not_played"
	if game.Status == "Jogado" {
		status = "played"
	}

	var steamAppID sql.NullInt64
	if game.SteamAppID != nil && *game.SteamAppID != "" {
		if val, err := strconv.ParseInt(*game.SteamAppID, 10, 64); err == nil {
			steamAppID = sql.NullInt64{Int64: val, Valid: true}
		}
	}

	source := "manual"
	if game.Platform == "steam" {
		source = "steam_sync"
	}

	_, err = tx.ExecContext(ctx, userGamesQuery,
		userID,
		catalogGameID,
		status,
		game.HoursPlayed,
		game.Rating,
		steamAppID,
		source,
	)
	if err != nil {
		return fmt.Errorf("failed to insert into user_games: %w", err)
	}

	return tx.Commit()
}
