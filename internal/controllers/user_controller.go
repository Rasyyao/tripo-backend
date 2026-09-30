package controllers

import (
	"errors"

	"tripo-backend/internal/repositories"
	"tripo-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	service services.UserService
}

func NewUserController(svc services.UserService) *UserController {
	return &UserController{service: svc}
}

func (ctrl *UserController) Create(c fiber.Ctx) error {
	var req createUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := ctrl.service.CreateUser(req.Name, req.Email)
	if err != nil {
		return mapServiceError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (ctrl *UserController) List(c fiber.Ctx) error {
	users, err := ctrl.service.ListUsers()
	if err != nil {
		return mapServiceError(err)
	}

	return c.JSON(users)
}

func (ctrl *UserController) Get(c fiber.Ctx) error {
	id := c.Params("id")

	user, err := ctrl.service.GetUser(id)
	if err != nil {
		return mapServiceError(err)
	}

	return c.JSON(user)
}

func (ctrl *UserController) Update(c fiber.Ctx) error {
	id := c.Params("id")

	var req updateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	user, err := ctrl.service.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		return mapServiceError(err)
	}

	return c.JSON(user)
}

func (ctrl *UserController) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.service.DeleteUser(id); err != nil {
		return mapServiceError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func mapServiceError(err error) error {
	switch {
	case errors.Is(err, services.ErrInvalidInput):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, repositories.ErrUserNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
