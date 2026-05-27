package postgres

import (
	"context"
	"database/sql"
	"errors"

	"your-next-game-backend/internal/identity/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert/Update users table
	userQuery := `
		INSERT INTO users (id, display_name, avatar_url, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (id) DO UPDATE
		SET display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
	`
	_, err = tx.ExecContext(ctx, userQuery, user.ID, user.Name, user.AvatarFull)
	if err != nil {
		return err
	}

	// 2. Insert/Update steam_accounts table
	steamQuery := `
		INSERT INTO steam_accounts (user_id, steam_id, persona_name, avatar_url, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (steam_id) DO UPDATE
		SET persona_name = EXCLUDED.persona_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
	`
	_, err = tx.ExecContext(ctx, steamQuery, user.ID, user.SteamID, user.Name, user.AvatarFull)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetBySteamID(ctx context.Context, steamID string) (*domain.User, error) {
	query := `
		SELECT u.id, s.steam_id, u.display_name, u.avatar_url
		FROM users u
		JOIN steam_accounts s ON u.id = s.user_id
		WHERE s.steam_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, steamID)

	var u domain.User
	err := row.Scan(&u.ID, &u.SteamID, &u.Name, &u.AvatarFull)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

