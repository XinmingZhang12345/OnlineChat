package main

import (
	"net/http"
	"sync"
	"time"

	restful "github.com/emicklei/go-restful/v3"
)

// Message represents a chat message.
type Message struct {
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatServer represents the chat server.
type ChatServer struct {
	users          map[string]string // Username to token mapping
	messages       []Message         // Message history
	userLastActive map[string]time.Time
	mu             sync.Mutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		users:          make(map[string]string),
		messages:       []Message{},
		userLastActive: make(map[string]time.Time),
	}
}

func main() {
	chatServer := NewChatServer()

	ws := new(restful.WebService)
	ws.Path("/chat").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// Register the chat service endpoints
	ws.Route(ws.GET("/login/{username}").To(chatServer.login))
	ws.Route(ws.POST("/logout").To(chatServer.logout))
	ws.Route(ws.POST("/message").To(chatServer.sendMessage))
	ws.Route(ws.GET("/messages").To(chatServer.getMessages))
	ws.Route(ws.GET("/users").To(chatServer.getUsers))

	restful.Add(ws)

	// Start a goroutine to check for inactive users and log them out
	go chatServer.checkInactiveUsers()

	// Start the server
	http.ListenAndServe(":8080", nil)
}

func (cs *ChatServer) checkInactiveUsers() {
}
