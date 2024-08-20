package main

import (
	"log"

	"backend/database"
	"backend/routes"
	"backend/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	app := fiber.New()

	routes.SetupRoutes(app, db)
	websocket.SetupWebSocketRoutes(app, rdb)

	log.Fatal(app.Listen(":4000"))
}
