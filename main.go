package main

import (
	"encoding/json"
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
func main() {
	
}