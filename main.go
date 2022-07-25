package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	websockethub "simple-websocket-echo-server/internal/websocket-hub"
)

func main() {
	app := fiber.New()

	handler := websockethub.NewHandler(app, 5, 5, 5, 50)
	go handler.Run()

	log.Fatal(app.Listen(":3000"))
}
