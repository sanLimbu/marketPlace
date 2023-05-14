package websocket

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	Messages   chan string
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		Messages:   make(chan string),
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
		s.Messages <- message
	}
}

func (s *Server) CloseWebSocket(ws *websocket.Conn) {
	s.unregister <- ws
	ws.Close()
}

func (s *Server) HandleConnections() {
	//	http.Handle("/ws", websocket.Handler(s.HandleWebSocket))
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
					//s.CloseWebSocket(conn) // Close the WebSocket connection
					conn.Close()
					delete(s.clients, conn)
				}
			}
		}
	}
}

func (s *Server) BroadcastMessages() chan<- []byte {
	fmt.Println("broadcast message kafka")
	return s.broadcast
}
