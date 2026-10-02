// Package middleware holds Fiber middleware shared across routes.
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"

	"tripo-backend/internal/auth"
	"tripo-backend/internal/response"
)

const userIDKey = "userID"

// Authenticate requires a valid "Authorization: Bearer <access token>" header
// and stores the user ID for downstream handlers (see UserID).
func Authenticate(jwt auth.TokenManager) fiber.Handler {
	return func(c fiber.Ctx) error {
		scheme, token, ok := strings.Cut(c.Get(fiber.HeaderAuthorization), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
			return unauthorized(c, "missing or malformed Authorization header")
		}

		claims, err := jwt.ParseAccess(strings.TrimSpace(token))
		if err != nil {
			return unauthorized(c, "invalid or expired access token")
		}

		c.Locals(userIDKey, claims.UserID)
		return c.Next()
	}
}

// UserID returns the authenticated user's ID set by Authenticate.
func UserID(c fiber.Ctx) string {
	id, _ := c.Locals(userIDKey).(string)
	return id
}

func unauthorized(c fiber.Ctx, message string) error {
	return response.Error(c, fiber.StatusUnauthorized, "unauthorized", message, nil)
}
