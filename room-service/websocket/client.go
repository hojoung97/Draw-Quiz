package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Name string
	Conn *websocket.Conn
	Hub  *Hub
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Fail to read client %s's message: %v\n", c.ID, err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}

		for client := range c.Hub.Clients {
			if c == client {
				continue
			}
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Fail in hub %d broadcast to user %s: %v\n", c.Hub.RoomID, client.ID, err)
			}
		}
	}
}
