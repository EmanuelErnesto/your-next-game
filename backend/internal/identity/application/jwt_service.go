package application

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateTokens(userID string) (accessToken string, refreshToken string, err error)
	VerifyToken(tokenString string) (userID string, err error)
}

type jwtService struct {
	secretKey []byte
}

func NewJWTService(secret string) TokenService {
	return &jwtService{
		secretKey: []byte(secret),
	}
}

func (s *jwtService) GenerateTokens(userID string) (string, string, error) {
	// Access Token (15 min)
	accessClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"typ": "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (7 days)
	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *jwtService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secretKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}
	return "", jwt.ErrTokenInvalidClaims
}
