package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"pocketapi/internal/config"
	"pocketapi/internal/database"
	"pocketapi/internal/routes"
)

import (
	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"pocketapi/internal/middleware"
)

func main() {

	config.LoadConfig()

	database.ConnectDatabase()

	app := fiber.New()
	app.Static("/uploads", "./uploads")
	app.Use(middleware.RequestID())

	app.Use(middleware.Logger())

	app.Use(recover.New())

	app.Use(fiberCors.New(middleware.CORS()))

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to PocketAPI ",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"status":   "OK",
			"database": "Connected",
		})
	})

	routes.Setup(app)

	port := fmt.Sprintf(":%s", config.AppConfig.Port)

	log.Printf("Server running on http://localhost%s\n", port)

	log.Fatal(app.Listen(port))
}
