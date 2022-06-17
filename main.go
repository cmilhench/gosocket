package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cmilhench/gosocket/uuid"

	"golang.org/x/net/websocket"
)

// ----------------------------------------------------------------------------

// Client stores the client's socket and information such as username
type Client struct {
	server    *SocketServer
	websocket *websocket.Conn
	send      chan string
	id        string
}

func (c *Client) Read() {
	for {
		var message string
		if err := websocket.Message.Receive(c.websocket, &message); err != nil {
			c.server.part <- c
			break // Error recieving, likely an EOF/disconnect, exit the loop
		}
		// Proccess message
		c.server.data <- message
	}
}

func (c *Client) Write() {
	for message := range c.send {
		websocket.Message.Send(c.websocket, message)
	}
}

func (c *Client) Send(message string) {
	select {
	case c.send <- message:
	default:
		c.server.part <- c
	}
}

func (c *Client) Close() {
	close(c.send)
	_ = c.websocket.Close()
}

// ----------------------------------------------------------------------------

// SocketServer holds all of the connected clients
type SocketServer struct {
	clients map[Client]struct{}
	data    chan string
	join    chan *Client
	part    chan *Client
}

// New creates a new SocketServer
func New() SocketServer {
	s := SocketServer{
		make(map[Client]struct{}),
		make(chan string),
		make(chan *Client),
		make(chan *Client),
	}
	go func() {
		for {
			select {
			case client := <-s.join:
				s.clients[*client] = struct{}{}
			case client := <-s.part:
				if _, ok := s.clients[*client]; ok {
					delete(s.clients, *client)
					client.Close()
				}
			case message := <-s.data:
				s.Broadcast(message)
			}
		}
	}()
	return s
}

func (s *SocketServer) Handler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		id, _ := uuid.New()
		client := &Client{s, ws, make(chan string), id}
		s.join <- client
		go client.Write()
		client.Read()
	})
}

func (s *SocketServer) Broadcast(message string) {
	for client := range s.clients {
		client.Send(message)
	}
}

// ----------------------------------------------------------------------------

func resolveAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf("localhost:%s", port)
}

func main() {
	socketServer := New()

	// go func() {
	// 	c := time.Tick(10 * time.Second)
	// 	for range c {
	// 		socketServer.Broadcast(time.Now().String())
	// 	}
	// }()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/socket", socketServer.Handler())

	server := &http.Server{Addr: resolveAddress()}
	log.Fatal(server.ListenAndServe())
}
