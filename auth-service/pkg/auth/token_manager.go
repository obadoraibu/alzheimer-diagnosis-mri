package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/obadoraibu/go-auth/internal/config"
)

type TokenManager struct {
	signingKey      string
	authTokenTTL    string
	refreshTokenTTL string
}

func NewTokenManager(cfg *config.AuthConfig) *TokenManager {
	return &TokenManager{
		signingKey:      cfg.SigningKey,
		authTokenTTL:    cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (tm *TokenManager) GenerateJWT(email, role string) (string, error) {

	duration, err := time.ParseDuration(tm.authTokenTTL)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(duration).Unix(),
		"role":  role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tm.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (tm *TokenManager) GenerateRefresh() string {
	id := uuid.New()
	return id.String()
}

func (tm *TokenManager) GetSigningKey() string {
	return tm.signingKey
}
