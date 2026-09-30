package routes

import (
	"gorm.io/gorm"

	"tripo-backend/internal/controllers"
	"tripo-backend/internal/repositories"
	"tripo-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
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
