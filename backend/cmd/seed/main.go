package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me-in-production-12345"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	userID := "11111111-1111-1111-1111-111111111111"
	
	// 1. Clean existing test data
	log.Println("Cleaning old seed data...")
	_, _ = db.Exec("DELETE FROM user_games WHERE user_id = $1", userID)
	_, _ = db.Exec("DELETE FROM users WHERE id = $1", userID)

	// 2. Insert test user
	log.Println("Inserting test user...")
	_, err = db.Exec(`
		INSERT INTO users (id, display_name, avatar_url)
		VALUES ($1, 'K6 Load Test User', 'https://example.com/avatar.jpg')
	`, userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	// 3. Insert mock games into catalog
	log.Println("Seeding games catalog...")
	for i := 1; i <= 200; i++ {
		gameID := fmt.Sprintf("%08d-%04d-%04d-%04d-%012d", i, i, i, i, i)
		providerGameID := fmt.Sprintf("%d", 1000+i)
		title := fmt.Sprintf("Awesome Game %d", i)
		coverURL := fmt.Sprintf("https://example.com/covers/%d.jpg", i)
		genre := "Action"
		if i%2 == 0 {
			genre = "RPG"
		} else if i%3 == 0 {
			genre = "Adventure"
		}
		
		_, err = db.Exec(`
			INSERT INTO games_catalog (id, provider, provider_game_id, title, cover_url, genre, description, developer, release_date)
			VALUES ($1, 'steam', $2, $3, $4, $5, 'This is a description for the game', 'Best Studio', '2022-01-01')
			ON CONFLICT (provider, provider_game_id) DO UPDATE SET title = EXCLUDED.title
		`, gameID, providerGameID, title, coverURL, genre)
		if err != nil {
			log.Fatalf("failed to insert game into catalog: %v", err)
		}

		// Associate first 50 games to user
		if i <= 50 {
			userGameID := fmt.Sprintf("%08d-f75a-4b62-bb34-%012d", i, i)
			status := "not_played"
			if i%2 == 0 {
				status = "played"
			}
			rating := float64(3 + (i % 3))
			hours := float64(i) * 2.5

			_, err = db.Exec(`
				INSERT INTO user_games (id, user_id, catalog_game_id, steam_app_id, status, hours_played, rating, source)
				VALUES ($1, $2, $3, $4, $5, $6, $7, 'steam_sync')
				ON CONFLICT (user_id, catalog_game_id) DO NOTHING
			`, userGameID, userID, gameID, 1000+i, status, hours, rating)
			if err != nil {
				log.Fatalf("failed to insert user game: %v", err)
			}
		}
	}

	// 4. Generate JWT Token
	log.Println("Generating JWT Access Token...")
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"typ": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("failed to sign JWT: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Printf("TEST_USER_ID: %s\n", userID)
	fmt.Printf("TEST_ACCESS_TOKEN: %s\n", tokenString)
	fmt.Println("==================================================")
}
