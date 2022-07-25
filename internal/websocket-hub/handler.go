package websocket_hub

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	client     map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan []byte
}

func NewHandler(app *fiber.App, clientSize, registerSize, unregisterSize, broadcastSize int) *Handler {
	client := make(map[*websocket.Conn]bool, clientSize)
	register := make(chan *websocket.Conn, registerSize)
	unregister := make(chan *websocket.Conn, unregisterSize)
	broadcast := make(chan []byte, broadcastSize)

	handler := &Handler{
		client:     client,
		register:   register,
		unregister: unregister,
		broadcast:  broadcast,
	}
	handler.registerRoute(app)

	return handler
}

func (H *Handler) Run() {
	for {
		select {
		case message := <-H.broadcast:
			for conn, _ := range H.client {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					fmt.Println("could not write to client: ", err)
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					conn.Close()
					delete(H.client, conn)
				}
			}
		case connection := <-H.register:
			H.client[connection] = true
			fmt.Println("new client registered")
		case connection := <-H.unregister:
			delete(H.client, connection)
			fmt.Println("client unregistered")
		}
	}
}

func (H *Handler) registerRoute(app *fiber.App) {
	app.Get("/ws", websocket.New(func(conn *websocket.Conn) {
		defer func() {
			H.unregister <- conn
			_ = conn.Close()
		}()
		H.register <- conn
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
					fmt.Println("unexpected error: ", err)
				}

				return
			}
			fmt.Println("message received: ", message)
			H.broadcast <- message
		}
	}))
}
