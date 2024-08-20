package routes

import (
	"backend/database"
	"backend/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, db *database.DB) {
	router := app.Group("/auth")

	router.Post("/login", auth.Login)
	router.Post("/register", auth.Register)
}
