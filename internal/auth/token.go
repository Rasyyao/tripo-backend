// Package auth issues and verifies the JWT access and refresh tokens.
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "tripo-backend"

// TokenType is stored in the "typ" claim so a refresh token can never be used
// as an access token (and vice versa).
type TokenType string

const (
	TypeAccess  TokenType = "access"
	TypeRefresh TokenType = "refresh"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type claims struct {
	jwt.RegisteredClaims
	Type TokenType `json:"typ"`
}

// Claims is the verified identity extracted from a token.
type Claims struct {
	UserID  string
	TokenID string // jti; for refresh tokens this is the refresh_tokens row ID
}

// IssuedToken is a freshly signed token and its metadata.
type IssuedToken struct {
	Token     string
	ID        string
	ExpiresAt time.Time
}

type TokenManager interface {
	GenerateAccess(userID string) (IssuedToken, error)
	GenerateRefresh(userID string) (IssuedToken, error)
	ParseAccess(token string) (Claims, error)
	ParseRefresh(token string) (Claims, error)
}

type jwtManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTManager(secret string, accessTTL, refreshTTL time.Duration) TokenManager {
	return &jwtManager{secret: []byte(secret), accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (m *jwtManager) GenerateAccess(userID string) (IssuedToken, error) {
	return m.generate(userID, TypeAccess, m.accessTTL)
}

func (m *jwtManager) GenerateRefresh(userID string) (IssuedToken, error) {
	return m.generate(userID, TypeRefresh, m.refreshTTL)
}

func (m *jwtManager) ParseAccess(token string) (Claims, error) {
	return m.parse(token, TypeAccess)
}

func (m *jwtManager) ParseRefresh(token string) (Claims, error) {
	return m.parse(token, TypeRefresh)
}

func (m *jwtManager) generate(userID string, typ TokenType, ttl time.Duration) (IssuedToken, error) {
	now := time.Now()
	id := uuid.NewString()
	expiresAt := now.Add(ttl)

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			ID:        id,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Type: typ,
	}).SignedString(m.secret)
	if err != nil {
		return IssuedToken{}, err
	}
	return IssuedToken{Token: signed, ID: id, ExpiresAt: expiresAt}, nil
}

func (m *jwtManager) parse(tokenStr string, want TokenType) (Claims, error) {
	c := &claims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		c,
		func(*jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid || c.Type != want || c.Subject == "" || c.ID == "" {
		return Claims{}, ErrInvalidToken
	}
	return Claims{UserID: c.Subject, TokenID: c.ID}, nil
}

// HashToken returns the hex SHA-256 of a token. Only this hash is persisted,
// so a database leak does not expose usable refresh tokens.
func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
