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

// NewChatServer creates a new chat server instance.
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

func (cs *ChatServer) login(req *restful.Request, resp *restful.Response) {
	// Extract the username from the request path
	username := req.PathParameter("username")

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Check if the username is already taken
	if _, exists := cs.users[username]; exists {
		resp.WriteHeader(http.StatusConflict)
		resp.WriteEntity("Username already in use")
		return
	}

	// Generate a token for the user
	token := generateToken()

	// Store the user and token
	cs.users[username] = token
	cs.userLastActive[username] = time.Now()

	resp.WriteHeader(http.StatusCreated)
	resp.WriteEntity(token)
}

func (cs *ChatServer) logout(req *restful.Request, resp *restful.Response) {

}

func (cs *ChatServer) sendMessage(req *restful.Request, resp *restful.Response) {

}

func (cs *ChatServer) getMessages(req *restful.Request, resp *restful.Response) {
}

func (cs *ChatServer) getUsers(req *restful.Request, resp *restful.Response) {

}

func (cs *ChatServer) checkInactiveUsers() {

}

func generateToken() string {

	return "sampletoken"
}
