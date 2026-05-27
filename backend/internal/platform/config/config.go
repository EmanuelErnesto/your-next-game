package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPPort             string
	AppVersion           string
	FrontendBaseURL      string
	FrontendSuccessURL   string
	SteamOpenIDRealm     string
	SteamOpenIDReturnURL string
	SteamAPIKey          string
	PostgresDSN          string
	LibraryRepository    string
	ReadHeaderTimeout    time.Duration
	OutboundHTTPTimeout  time.Duration
	JWTSecret            string
}

func loadEnv(filenames ...string) {
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
				val = val[1 : len(val)-1]
			}

			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}

func Load() Config {
	loadEnv(".env", "../.env", "your-next-game-backend/.env")

	frontendBaseURL := getEnv("FRONTEND_BASE_URL", "http://localhost:3001")
	frontendSuccessURL := getEnv("FRONTEND_SUCCESS_URL", frontendBaseURL+"/")
	
	httpPort := getEnv("HTTP_PORT", "8080")
	steamOpenIDRealm := getEnv("STEAM_OPENID_REALM", "http://localhost:"+httpPort)
	steamOpenIDReturnURL := getEnv("STEAM_OPENID_RETURN_URL", steamOpenIDRealm+"/v1/auth/steam/callback")


	return Config{
		HTTPPort:             httpPort,
		AppVersion:           getEnv("APP_VERSION", "dev"),
		FrontendBaseURL:      frontendBaseURL,
		FrontendSuccessURL:   frontendSuccessURL,
		SteamOpenIDRealm:     steamOpenIDRealm,
		SteamOpenIDReturnURL: steamOpenIDReturnURL,
		SteamAPIKey:          getEnv("STEAM_API_KEY", ""),
		PostgresDSN:          getEnv("POSTGRES_DSN", ""),
		LibraryRepository:    getEnv("LIBRARY_REPOSITORY", "postgres"),
		ReadHeaderTimeout:    3 * time.Second,
		OutboundHTTPTimeout:  8 * time.Second,
		JWTSecret:            getEnv("JWT_SECRET", "change-me-in-production-12345"),
	}
}

func (c Config) Validate() error {
	if c.LibraryRepository != "postgres" {
		return fmt.Errorf("unsupported LIBRARY_REPOSITORY %q", c.LibraryRepository)
	}
	if c.PostgresDSN == "" {
		return fmt.Errorf("POSTGRES_DSN is required when LIBRARY_REPOSITORY=postgres")
	}
	return nil
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
