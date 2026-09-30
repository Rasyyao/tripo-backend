package routes

import (
	"tripo-backend/internal/controllers"
	"tripo-backend/internal/repositories"
	"tripo-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

// Setup wires repositories, services and controllers, then registers routes.
func Setup(app *fiber.App) {
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	api := app.Group("/api/v1")

	users := api.Group("/users")
	users.Post("/", userController.Create)
	users.Get("/", userController.List)
	users.Get("/:id", userController.Get)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)
}
