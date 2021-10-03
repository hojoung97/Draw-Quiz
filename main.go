package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hojoung97/Draw-Quiz/pkg/websocket"
)

var hubs map[int]*websocket.Hub

func setupWS(roomID int) {
	if _, ok := hubs[roomID]; !ok {
		hubs[roomID] = websocket.NewHub()
		go hubs[roomID].Start()
	}
}

func getRoomID(path string) (int, error) {
	parsed := strings.Split(path, "/")
	roomID, err := strconv.Atoi(parsed[len(parsed)-1])
	if err != nil {
		return 0, err
	}
	return roomID, nil
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Websocket Endpoint Hit: %s in %s\n", r.Host, r.URL.Path)

	// Upgrade the connection to a WebSocket connection
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	roomID, err := getRoomID(r.URL.Path)
	if err != nil {
		log.Fatal(err)
	}
	hub := hubs[roomID]

	client := &websocket.Client{
		Conn: conn,
		Hub:  hub,
	}

	hub.Register <- client
	// listen indefinitely for new messages on our websocket conn
	client.Read()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleRoom(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	roomID, err := strconv.Atoi(query.Get("roomID"))
	if err != nil {
		log.Fatal(err)
	}

	//if hub, ok := hubs[roomID]; ok && hub.Capacity == 2 {
	//}
	setupWS(roomID)
	http.ServeFile(w, r, "draw_app.html")
}

func main() {
	hubs = make(map[int]*websocket.Hub)

	webServerMux := mux.NewRouter()

	webServerMux.HandleFunc("/", handleHome)
	webServerMux.HandleFunc("/room", handleRoom).Methods("GET")
	webServerMux.HandleFunc("/room/{roomID:[0-9]+}", handleWS)

	fmt.Println("Draw App Web Server Starting...")
	http.ListenAndServe(":8080", webServerMux)
}
