package services

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"tripo-backend/internal/apperror"
	"tripo-backend/internal/auth"
	"tripo-backend/internal/dto"
	"tripo-backend/internal/models"
	"tripo-backend/internal/repositories"
)

const tokenTypeBearer = "Bearer"

var (
	errInvalidCredentials = apperror.New(apperror.ErrUnauthorized, "invalid email or password")
	errInvalidRefresh     = apperror.New(apperror.ErrUnauthorized, "invalid or expired refresh token")

	// dummyHash lets Login spend the same bcrypt time for unknown emails,
	// so response timing does not reveal which emails are registered.
	dummyHash, _ = bcrypt.GenerateFromPassword([]byte("tripo-dummy-password"), bcrypt.DefaultCost)
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
	Refresh(req dto.RefreshRequest) (*dto.TokenResponse, error)
	Logout(req dto.RefreshRequest) (*dto.MessageResponse, error)
}

type authService struct {
	users  repositories.UserRepository
	tokens repositories.RefreshTokenRepository
	jwt    auth.TokenManager
}

func NewAuthService(users repositories.UserRepository, tokens repositories.RefreshTokenRepository, jwt auth.TokenManager) AuthService {
	return &authService{users: users, tokens: tokens, jwt: jwt}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := normalizeEmail(req.Email)

	if err := collect(validateName(name), validateEmail(email), validatePassword(req.Password)); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.users.Create(&models.User{
		DisplayName:  &name,
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	return s.authResponse(user)
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	if err := collect(validateRequired("email", email), validateRequired("password", req.Password)); err != nil {
		return nil, err
	}

	user, err := s.users.FindByEmail(email)
	if err != nil {
		if isNotFound(err) {
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
			return nil, errInvalidCredentials
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, errInvalidCredentials
	}
	return s.authResponse(user)
}

// Refresh rotates the refresh token: the presented one is revoked and a new
// pair is issued. Presenting an already-revoked token means it was replayed
// (likely stolen), so every session of that user is revoked.
func (s *authService) Refresh(req dto.RefreshRequest) (*dto.TokenResponse, error) {
	if err := collect(validateRequired("refresh_token", req.RefreshToken)); err != nil {
		return nil, err
	}

	claims, err := s.jwt.ParseRefresh(req.RefreshToken)
	if err != nil {
		return nil, errInvalidRefresh
	}

	stored, err := s.tokens.FindByHash(auth.HashToken(req.RefreshToken))
	if err != nil {
		if isNotFound(err) {
			return nil, errInvalidRefresh
		}
		return nil, err
	}
	if stored.ID != claims.TokenID || stored.UserID != claims.UserID {
		return nil, errInvalidRefresh
	}

	if stored.RevokedAt != nil {
		if err := s.tokens.RevokeAllForUser(stored.UserID); err != nil {
			return nil, err
		}
		return nil, errInvalidRefresh
	}
	if time.Now().After(stored.ExpiresAt) {
		return nil, errInvalidRefresh
	}

	revoked, err := s.tokens.Revoke(stored.ID)
	if err != nil {
		return nil, err
	}
	if !revoked { // lost a race with a concurrent refresh using the same token
		return nil, errInvalidRefresh
	}

	user, err := s.users.FindByID(stored.UserID)
	if err != nil {
		if isNotFound(err) {
			return nil, errInvalidRefresh
		}
		return nil, err
	}
	return s.issueTokens(user.ID)
}

// Logout revokes the presented refresh token. It is idempotent: unknown,
// expired or already-revoked tokens are not an error.
func (s *authService) Logout(req dto.RefreshRequest) (*dto.MessageResponse, error) {
	done := dto.NewMessageResponse("logged out successfully")

	if err := collect(validateRequired("refresh_token", req.RefreshToken)); err != nil {
		return nil, err
	}
	if _, err := s.jwt.ParseRefresh(req.RefreshToken); err != nil {
		return done, nil
	}

	stored, err := s.tokens.FindByHash(auth.HashToken(req.RefreshToken))
	if err != nil {
		if isNotFound(err) {
			return done, nil
		}
		return nil, err
	}
	if _, err := s.tokens.Revoke(stored.ID); err != nil {
		return nil, err
	}
	return done, nil
}

func (s *authService) authResponse(user *models.User) (*dto.AuthResponse, error) {
	tokens, err := s.issueTokens(user.ID)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{User: dto.NewUserResponse(user), TokenResponse: *tokens}, nil
}

// issueTokens signs a new access/refresh pair and stores the refresh token's hash.
func (s *authService) issueTokens(userID string) (*dto.TokenResponse, error) {
	access, err := s.jwt.GenerateAccess(userID)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwt.GenerateRefresh(userID)
	if err != nil {
		return nil, err
	}

	err = s.tokens.Create(&models.RefreshToken{
		ID:        refresh.ID,
		UserID:    userID,
		TokenHash: auth.HashToken(refresh.Token),
		ExpiresAt: refresh.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
		TokenType:    tokenTypeBearer,
		ExpiresIn:    int(time.Until(access.ExpiresAt).Seconds()),
	}, nil
}
