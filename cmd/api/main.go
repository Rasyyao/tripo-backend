package main

import (
	"log"

	"tripo-backend/internal/config"
	"tripo-backend/internal/database"
	"tripo-backend/internal/routes"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	routes.Setup(app, db)

	log.Printf("server listening on :%s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
