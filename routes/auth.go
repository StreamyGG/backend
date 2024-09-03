package routes

import (
	"backend/database"
	"backend/handlers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, db *database.DB) {
	router := app.Group("/auth")

	router.Post("/login", func(c *fiber.Ctx) error {
		return auth.Login(c, db)
	})
	router.Post("/register", func(c *fiber.Ctx) error {
		return auth.Register(c, db)
	})
}
