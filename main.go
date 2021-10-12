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

	vars := mux.Vars(r)

	roomID, err := strconv.Atoi(vars["roomID"])
	if err != nil {
		log.Printf("Fail to convert the roomID to int (handleWS): %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &websocket.Client{
		Name: vars["userName"],
		Conn: conn,
		Hub:  hubs[roomID],
	}

	hubs[roomID].Register <- client
	// listen indefinitely for new messages on our websocket conn
	client.Read()
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
	http.ServeFile(w, r, "static/draw_app.html")
}

func main() {
	hubs = make(map[int]*websocket.Hub)

	webServerMux := mux.NewRouter()

	webServerMux.HandleFunc("/room", handleRoom).Methods("GET")
	webServerMux.HandleFunc("/room/{roomID:[0-9]+}/{userName}", handleWS)
	webServerMux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	// TODO: Make port as a configurable parameter
	log.Printf("Draw App Web Server Listening on localhost%s\n", ":8080")
	http.ListenAndServe(":8080", webServerMux)
}
