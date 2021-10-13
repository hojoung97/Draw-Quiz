FROM golang:alpine AS builder
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go get ./...

COPY room-service/main.go ./
ADD room-service/websocket ./websocket
RUN go build -o room-service

FROM alpine
COPY --from=builder /build/room-service ./
ENTRYPOINT [""]
