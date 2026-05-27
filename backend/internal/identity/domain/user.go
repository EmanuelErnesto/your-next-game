package domain

import "context"

type User struct {
	ID         string
	SteamID    string
	Name       string
	AvatarFull string
}

type UserRepository interface {
	Save(ctx context.Context, user User) error
	GetBySteamID(ctx context.Context, steamID string) (*User, error)
}
