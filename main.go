package main

import (
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	"time"
	"fmt"
)

type loginReqest struct {
	Username string 'json:username'
}

type authToken struct {
	Token string 'json:token'
}

type userStatus struct {
	Username string 'json:username',
	Online string 'json:online'
}

type message struct {
	Message string 'json:message'
}

type messageResponse struct {
	Id int 'json:id',
	Message string 'json:message',
	Author string 'json:author'
}

func generateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func handleLogin(request *restful.request) {
}
func handlSendMessage(request *restful.request){

}
func handleGetMessage(request *restful.request){

}
func handleLogout(request *restful.request){
	
}
func main() {
	
}