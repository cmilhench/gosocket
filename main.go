package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

// Client stores the client's socket and information such as username
type Client struct {
	websocket *websocket.Conn
	clientIP  string
}

// SocketServer holds all of the connected clients
type SocketServer struct {
	clients map[Client]struct{}
}

// New creates a new SocketServer
func New() http.Handler {
	s := SocketServer{make(map[Client]struct{})}
	return websocket.Handler(func(ws *websocket.Conn) {
		var message string
		defer func() { _ = ws.Close() }()
		client := Client{ws, ws.Request().RemoteAddr}
		s.clients[client] = struct{}{}
		for {
			if err := websocket.Message.Receive(client.websocket, &message); err != nil {
				delete(s.clients, client)
				break
			}
			// Procccess message by echoing to all.
			for client := range s.clients {
				if err := websocket.Message.Send(client.websocket, client.clientIP+message); err != nil {
					log.Printf("Can't send %v\n", err)
				}
			}
		}
	})
}

func resolveAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf("localhost:%s", port)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/socket", New())

	server := &http.Server{Addr: resolveAddress()}
	log.Println(server.ListenAndServe())
}
