package controllers

import (
	"tripo-backend/internal/dto"
	"tripo-backend/internal/response"
	"tripo-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(svc services.AuthService) *AuthController {
	return &AuthController{service: svc}
}

func (ctrl *AuthController) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	res, err := ctrl.service.Register(req)
	if err != nil {
		return respondError(c, err)
	}
	return response.Created(c, res)
}

func (ctrl *AuthController) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	res, err := ctrl.service.Login(req)
	if err != nil {
		return respondError(c, err)
	}
	return response.OK(c, res)
}

func (ctrl *AuthController) Refresh(c fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	res, err := ctrl.service.Refresh(req)
	if err != nil {
		return respondError(c, err)
	}
	return response.OK(c, res)
}

func (ctrl *AuthController) Logout(c fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if err := ctrl.service.Logout(req); err != nil {
		return respondError(c, err)
	}
	return response.NoContent(c)
}
