package main

import (
	"log"

	"tripo-backend/internal/config"
	"tripo-backend/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()

	app := fiber.New()

	routes.Setup(app)

	log.Printf("server listening on :%s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
