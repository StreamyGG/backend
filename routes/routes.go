package routes

import (
	"backend/database"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *database.DB) {

	SetupAuthRoutes(app, db)
}
