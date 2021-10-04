package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hojoung97/Draw-Quiz/pkg/websocket"
)

var hubs map[int]*websocket.Hub

func setupWS(roomID int) {
	hubs[roomID] = websocket.NewHub(roomID)
	go hubs[roomID].Start()
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Websocket Endpoint Hit: %s%s\n", r.Host, r.URL.Path)

	// Upgrade the connection to a WebSocket connection
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Printf("Fail to upgrade the connection (handleWS): %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	roomID, err := strconv.Atoi(mux.Vars(r)["roomID"])
	if err != nil {
		log.Printf("Fail to convert the roomID to int (handleWS): %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Hub:  hubs[roomID],
	}

	hubs[roomID].Register <- client
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
		log.Printf("Fail to convert the roomID to int (handleRoom): %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: set maximum capacity per room
	if hub, ok := hubs[roomID]; !ok || hub == nil {
		setupWS(roomID)
	}
	http.ServeFile(w, r, "draw_app.html")
}

func main() {
	hubs = make(map[int]*websocket.Hub)

	webServerMux := mux.NewRouter()

	webServerMux.HandleFunc("/", handleHome)
	webServerMux.HandleFunc("/room", handleRoom).Methods("GET")
	webServerMux.HandleFunc("/room/{roomID:[0-9]+}", handleWS)

	// TODO: Make port as a configurable parameter
	log.Printf("Draw App Web Server Listening on localhost%s\n", ":8080")
	http.ListenAndServe(":8080", webServerMux)
}
