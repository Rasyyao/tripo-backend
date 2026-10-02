package routes

import (
	"gorm.io/gorm"

	"tripo-backend/internal/auth"
	"tripo-backend/internal/config"
	"tripo-backend/internal/controllers"
	"tripo-backend/internal/middleware"
	"tripo-backend/internal/repositories"
	"tripo-backend/internal/services"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	userRepo := repositories.NewUserRepository(db)
	refreshRepo := repositories.NewRefreshTokenRepository(db)

	authController := controllers.NewAuthController(services.NewAuthService(userRepo, refreshRepo, jwtManager))
	userController := controllers.NewUserController(services.NewUserService(userRepo))

	api := app.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.Post("/register", authController.Register)
	authRoutes.Post("/login", authController.Login)
	authRoutes.Post("/refresh", authController.Refresh)
	authRoutes.Post("/logout", authController.Logout)

	users := api.Group("/users", middleware.Authenticate(jwtManager))
	users.Get("/", userController.List)
	users.Get("/:id", userController.Get)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)
}
