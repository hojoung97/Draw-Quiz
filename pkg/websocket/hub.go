package websocket

import (
	"fmt"
	"log"
)

type Hub struct {
	RoomID     int
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub(roomID int) *Hub {
	return &Hub{
		RoomID:     roomID,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (hub *Hub) Start() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true
			log.Printf("New Client %s joined room %d (size=%d)\n", client.ID, hub.RoomID, len(hub.Clients))
			for client := range hub.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("New User %s Joined", client.ID)})
			}
		case client := <-hub.Unregister:
			delete(hub.Clients, client)
			log.Printf("New Client %s left room %d (size=%d)\n", client.ID, hub.RoomID, len(hub.Clients))
			for client := range hub.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: fmt.Sprintf("User %s Disconnected", client.ID)})
			}
		case message := <-hub.Broadcast:
			// log.Printf("Sending message to all clients in hub %d\n", hub.RoomID)
			for client := range hub.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Printf("Fail in hub %d broadcast to user %s: %v\n", hub.RoomID, client.ID, err)
				}
			}
		}
	}
}
