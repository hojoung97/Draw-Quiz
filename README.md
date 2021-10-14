# Draw-Quiz
Draw Quiz Web App

## Requirements for Running the App

### Running Locally
`go 1.16` or higher

### Running on Docker
[`docker`](https://docs.docker.com/get-docker/)

## Frontend
Drawing quiz game as a web application 

### Made with
- HTML, CSS and Vanilla JS
- Golang web server app with Gorilla Mux

### Start Frontend Service
```
cd webserver
```
```
WEBSERVER_PORT=8000 go run main.go
```

## Game Room Service
- Verifies and handles users' room entry requests
- Opens websocket connections for each client in the room
- Clients communicate through the open websocket connection

### Made with
- Golang and Gorilla Websocket

### Start Game Room Service
```
cd ws-service
```
```
WEBSOCKET_PORT=8050 go run main.go
```
NOTE: currently does not suppot other port. Please use port 8050 for now.

## Docker
- Containerize each services for ease
- Apply basic concept of mircoservice architecture

To run the application with Docker:
```
docker-compose up
```

Once the containers are up and running, open `localhost:8000` on your browser

NOTE: if you changed the port for the webserver, use that port ex. `localhost:${WEBSERVER_PORT}`)

`Dockerfile`s for frontend and game room service are under `webserver` and `ws-service` directory, respectively. 

## WIP
This is a very rough version but please enjoy!
