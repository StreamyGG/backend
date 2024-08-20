package websocket

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func SetupWebSocketRoutes(app *fiber.App, rdb *redis.Client) {
	app.Get("/ws", HandleChat(rdb))

}
