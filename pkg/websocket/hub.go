package websocket

import "fmt"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
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
			fmt.Println("Size of Connection Hub: ", len(hub.Clients))
			for client := range hub.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined"})
			}
		case client := <-hub.Unregister:
			delete(hub.Clients, client)
			fmt.Println("Size of Connection Hub: ", len(hub.Clients))
			for client := range hub.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected"})
			}
		case message := <-hub.Broadcast:
			fmt.Println("Sending message to all clients in hub")
			for client := range hub.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
