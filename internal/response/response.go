// Package response defines the JSON envelope shared by every endpoint.
//
//	success: {"success": true,  "data": ..., "meta": {...}}
//	failure: {"success": false, "error": {"code": "...", "message": "...", "details": [...]}}
package response

import "github.com/gofiber/fiber/v3"

type Meta struct {
	Total int `json:"total"`
}

type Success struct {
	Success bool  `json:"success"`
	Data    any   `json:"data"`
	Meta    *Meta `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Failure struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

func OK(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Success{Success: true, Data: data})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Success{Success: true, Data: data})
}

// List wraps a collection; total is the number of items in the collection.
func List(c fiber.Ctx, data any, total int) error {
	return c.Status(fiber.StatusOK).JSON(Success{Success: true, Data: data, Meta: &Meta{Total: total}})
}

func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func Error(c fiber.Ctx, status int, code, message string, details any) error {
	return c.Status(status).JSON(Failure{
		Success: false,
		Error:   ErrorBody{Code: code, Message: message, Details: details},
	})
}
