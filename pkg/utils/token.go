package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/Stenoliv/didlydoodash_api/internal/config"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenParams struct {
	UserID     string
	RememberMe bool
}

// Generate a new access token
func GenerateAccessToken(cfg *config.EnvConfig, params TokenParams) (string, error) {
	lifespan := cfg.TokenAccessTTL

	// Generate values
	jti := gonanoid.Must()
	exp := time.Now().Add(lifespan)

	claims := jwt.MapClaims{}
	claims["jti"] = jti
	claims["sub"] = params.UserID
	claims["iss"] = "didlydoodash_api"
	claims["aud"] = "didlydoodash_frontend"
	claims["exp"] = exp.Unix()
	claims["type"] = AccessToken
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.TokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign %s token: %w", AccessToken, err)
	}
	return t, nil
}

// Generate a new refresh token
func GenerateRefreshToken(cfg *config.EnvConfig, params TokenParams) (string, error) {
	var lifespan time.Duration
	if params.RememberMe {
		lifespan = cfg.TokenRefreshRememberTTL
	} else {
		lifespan = cfg.TokenRefreshTTL
	}

	// Generate values
	jti := gonanoid.Must()
	exp := time.Now().Add(lifespan)

	claims := jwt.MapClaims{}
	claims["jti"] = jti
	claims["sub"] = params.UserID
	claims["iss"] = "didlydoodash_api"
	claims["aud"] = "didlydoodash_frontend"
	claims["exp"] = exp.Unix()
	claims["remember"] = params.RememberMe
	claims["type"] = RefreshToken
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg.TokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign %s token: %w", RefreshToken, err)
	}
	return t, nil
}

// Exctract access token from cookie or Authorization header
func ExtractToken(c *gin.Context) string {
	if access_token, err := c.Request.Cookie("token"); err == nil {
		return access_token.Value
	}
	if bearerToken := c.Request.Header.Get("Authorization"); bearerToken != "" && len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	if queryToken := c.Query("token"); queryToken != "" {
		return queryToken
	}
	return ""
}

func ValidateToken(cfg *config.EnvConfig, tokenStr string, expectedType TokenType) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(cfg.TokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims["type"] != string(expectedType) {
		return nil, fmt.Errorf("token type mismatch")
	}

	return &claims, nil
}
