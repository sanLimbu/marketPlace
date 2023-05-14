package websocket

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (s *Server) HandleWebSocket(ws *websocket.Conn) {
	s.register <- ws
	defer func() {
		s.unregister <- ws
		ws.Close()
	}()

	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			log.Printf("WebSocket receive error: %v", err)
			break
		}

		// Process the received message
		log.Printf("Received message from WebSocket client: %s", message)
		s.broadcast <- []byte(message)
	}
}

func (s *Server) HandleConnections() {
	for {
		select {
		case conn := <-s.register:
			s.clients[conn] = true
		case conn := <-s.unregister:
			if _, ok := s.clients[conn]; ok {
				delete(s.clients, conn)
			}
		case message := <-s.broadcast:
			for conn := range s.clients {
				err := websocket.Message.Send(conn, string(message))
				if err != nil {
					log.Printf("WebSocket send error: %v", err)
					conn.Close()
					delete(s.clients, conn)
				}
			}
		}
	}
}

func (s *Server) BroadcastMessages() chan<- []byte {
	fmt.Println("broadcast: ", s.broadcast)
	return s.broadcast
}
