package main

import (
	"fmt"
	"net/http"

	"github.com/hojoung97/Draw-Quiz/pkg/websocket"
)

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket Endpoint Hit: ", r.Host)

	// Upgrade the connection to a WebSocket connection
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	// listen indefinitely for new messages on our websocket conn
	client.Read()
}

func setupRoutes() *http.ServeMux {
	webServerMux := http.NewServeMux()
	pool := websocket.NewPool()
	go pool.Start()
	webServerMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
	return webServerMux
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleRoom(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "draw_app.html")
}

func main() {
	webServerMux := setupRoutes()
	webServerMux.HandleFunc("/", handleHome)
	webServerMux.HandleFunc("/room", handleRoom)

	fmt.Println("Draw App Web Server Starting...")
	http.ListenAndServe(":8080", webServerMux)
}
