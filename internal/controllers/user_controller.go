package controllers

import (
	"tripo-backend/internal/dto"
	"tripo-backend/internal/response"
	"tripo-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	service services.UserService
}

func NewUserController(svc services.UserService) *UserController {
	return &UserController{service: svc}
}

func (ctrl *UserController) List(c fiber.Ctx) error {
	users, err := ctrl.service.ListUsers()
	if err != nil {
		return respondError(c, err)
	}

	return response.List(c, users, len(users))
}

func (ctrl *UserController) Get(c fiber.Ctx) error {
	id := c.Params("id")

	user, err := ctrl.service.GetUser(id)
	if err != nil {
		return respondError(c, err)
	}

	return response.OK(c, user)
}

func (ctrl *UserController) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	user, err := ctrl.service.UpdateUser(id, req)
	if err != nil {
		return respondError(c, err)
	}

	return response.OK(c, user)
}

func (ctrl *UserController) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.service.DeleteUser(id); err != nil {
		return respondError(c, err)
	}

	return response.NoContent(c)
}
