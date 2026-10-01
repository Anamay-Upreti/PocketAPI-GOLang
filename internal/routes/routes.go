package routes

import (
	"github.com/gofiber/fiber/v2"

	"pocketapi/internal/handlers"
	"pocketapi/internal/middleware"
)

func Setup(app *fiber.App) {

	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	noteHandler := handlers.NewNoteHandler()

	api := app.Group("/api")

	// =========================
	// AUTH ROUTES
	// =========================

	auth := api.Group("/auth")

	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// =========================
	// USER ROUTES
	// =========================

	users := api.Group("/users")

	users.Use(middleware.Protected())

	users.Get("/me", userHandler.GetProfile)
	users.Post("/avatar", userHandler.UploadAvatar)

	// Logout from all devices
	users.Post("/logout-all", authHandler.LogoutAll)

	// =========================
	// NOTE ROUTES
	// =========================

	notes := api.Group("/notes")

	notes.Use(middleware.Protected())

	notes.Post("/", noteHandler.Create)
	notes.Get("/", noteHandler.GetAll)
	notes.Get("/:id", noteHandler.GetOne)
	notes.Put("/:id", noteHandler.Update)
	notes.Delete("/:id", noteHandler.Delete)
}