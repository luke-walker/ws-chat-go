package main

import (
	"fmt"
    "io"
	"net/http"

	"golang.org/x/net/websocket"
)

const ADDR = "127.0.0.1:5050"

type Server struct {
	clients map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) ChatServer(ws *websocket.Conn) {
    // should have client send initial data containing username etc
	s.clients[ws] = true
	
    buf := make([]byte, 1024)
    for {
        n, err := ws.Read(buf)
        if n == 0 {
            fmt.Printf("Timeout: Disconnecting %p\n", ws)
            delete(s.clients, ws)
            break
        }
        if err != nil && err != io.EOF {
            fmt.Println("Error:", err)
        }

        // can optimize w/ threads
        for client := range s.clients {
            if client == ws {
                continue
            }
            
            _, err := client.Write(buf[:n])
            if err != nil {
                fmt.Println("Error:", err)
            }
        }
    }
}

func main() {
	server := NewServer()

	http.Handle("/", websocket.Handler(server.ChatServer))
	
	fmt.Println("Server is running on", ADDR)
	err := http.ListenAndServe(ADDR, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
