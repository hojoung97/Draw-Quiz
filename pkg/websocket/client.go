package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
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
		c.Hub.Broadcast <- message
		//log.Printf("Message Received: %+v\n", message)
	}
}
