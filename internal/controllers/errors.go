package controllers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"

	"tripo-backend/internal/apperror"
	"tripo-backend/internal/response"
)

func badRequest(c fiber.Ctx, message string) error {
	return response.Error(c, fiber.StatusBadRequest, "bad_request", message, nil)
}

// respondError translates a service/repository error into the standard error
// envelope. Unknown errors are logged and hidden behind a generic message.
func respondError(c fiber.Ctx, err error) error {
	var validationErr *apperror.ValidationError

	switch {
	case errors.As(err, &validationErr):
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation_failed",
			"one or more fields are invalid", validationErr.Fields)
	case errors.Is(err, apperror.ErrNotFound):
		return response.Error(c, fiber.StatusNotFound, "not_found", err.Error(), nil)
	case errors.Is(err, apperror.ErrConflict):
		return response.Error(c, fiber.StatusConflict, "conflict", err.Error(), nil)
	default:
		log.Printf("internal error: %v", err)
		return response.Error(c, fiber.StatusInternalServerError, "internal_error", "internal server error", nil)
	}
}
