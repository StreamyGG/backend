package main

import (
	"log"

	"backend/database"
	"backend/routes"
	"backend/websocket"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

var (
	redisAddr string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err);
	}
	redisAddr = os.GetEnv("REDIS_ADDR")
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	app := fiber.New()

	routes.SetupRoutes(app, db)
	websocket.SetupWebSocketRoutes(app, rdb)

	log.Fatal(app.Listen(":4000"))
}
