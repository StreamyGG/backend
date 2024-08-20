package websocket

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var ctx = context.Background()

func HandleChat(rdb *redis.Client) func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		pubsub := rdb.Subscribe(ctx, "chat")
		defer pubsub.Close()

		ch := pubsub.Channel()

		c.SetReadDeadline(time.Now().Add(30 * time.Second))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(30 * time.Second))
			return nil
		})

		go func() {
			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					log.Println("read error:", err)
					break
				}
			}
		}()

		for msg := range ch {
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Println("write error:", err)
				break
			}
		}
	})
}
