package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type TokenManager interface {
	Generate(userID string) (string, error)
	Parse(token string) (userID string, err error)
}

type jwtManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) TokenManager {
	return &jwtManager{secret: []byte(secret), ttl: ttl}
}

func (m *jwtManager) Generate(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(m.secret)
}

func (m *jwtManager) Parse(tokenStr string) (string, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(*jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || !token.Valid || claims.Subject == "" {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}
