package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/google/uuid"
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
	username := req.QueryParameter("username")

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Check if the user exists
	if _, exists := cs.users[username]; exists {
		delete(cs.users, username)
		delete(cs.userLastActive, username)
		resp.WriteHeader(http.StatusOK)
		resp.WriteEntity("Logout successful")
	} else {
		resp.WriteHeader(http.StatusNotFound)
		resp.WriteEntity("User not found")
	}
}

func (cs *ChatServer) sendMessage(req *restful.Request, resp *restful.Response) {
	message := Message{}
	err := req.ReadEntity(&message)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		resp.WriteEntity("Invalid message format")
		return
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Check if the sender exists
	if _, exists := cs.users[message.Sender]; !exists {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.WriteEntity("User not logged in")
		return
	}

	message.Timestamp = time.Now()
	cs.messages = append(cs.messages, message)
	cs.userLastActive[message.Sender] = time.Now()

	resp.WriteHeader(http.StatusCreated)
	resp.WriteEntity("Message sent successfully")
}

func (cs *ChatServer) getMessages(req *restful.Request, resp *restful.Response) {
	username := req.QueryParameter("username")

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Check if the user exists
	if _, exists := cs.users[username]; !exists {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.WriteEntity("User not logged in")
		return
	}

	// Filter messages sent after the user's last activity
	userLastActive := cs.userLastActive[username]
	filteredMessages := []Message{}
	for _, message := range cs.messages {
		if message.Timestamp.After(userLastActive) {
			filteredMessages = append(filteredMessages, message)
		}
	}

	resp.WriteHeader(http.StatusOK)
	resp.WriteEntity(filteredMessages)
}

func (cs *ChatServer) getUsers(req *restful.Request, resp *restful.Response) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Get the list of currently logged-in users
	userList := make([]string, 0, len(cs.users))
	for user := range cs.users {
		userList = append(userList, user)
	}

	resp.WriteHeader(http.StatusOK)
	resp.WriteEntity(userList)
}

func (cs *ChatServer) checkInactiveUsers() {
	for {
		time.Sleep(time.Minute)

		cs.mu.Lock()
		currentTime := time.Now()
		for username, lastActiveTime := range cs.userLastActive {
			if currentTime.Sub(lastActiveTime) > (5 * time.Minute) { // Change the timeout duration as needed
				delete(cs.users, username)
				delete(cs.userLastActive, username)
				fmt.Printf("User %s logged out due to inactivity\n", username)
			}
		}
		cs.mu.Unlock()
	}
}

func generateToken() string {

	return uuid.New().String()
}
