package websocket

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Define our Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return conn, err
	}
	return conn, nil
}

func Reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("pkg/websocket Reader: %v\n", err)
			return
		}
		fmt.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Printf("pkg/websocket Reader: %v\n", err)
			return
		}
	}
}

func Writer(conn *websocket.Conn) {
	for {
		messageType, r, err := conn.NextReader()
		if err != nil {
			log.Printf("pkg/websocket Writer: %v\n", err)
			return
		}
		w, err := conn.NextWriter(messageType)
		if err != nil {
			log.Printf("pkg/websocket Writer: %v\n", err)
			return
		}
		if _, err := io.Copy(w, r); err != nil {
			log.Printf("pkg/websocket Writer: %v\n", err)
			return
		}
		if err := w.Close(); err != nil {
			log.Printf("pkg/websocket Writer: %v\n", err)
			return
		}
	}
}
