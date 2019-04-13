package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type Nothing bool

type ChatServer struct {
	users []string
	messageQueue map[string][]string
	mutex sync.Mutex
} 
type Message struct {
	Nickname   string
	Text string
}

func (chat *ChatServer) CreateUser(nickname string, reply *string) error {
	chat.mutex.Lock()
	defer chat.mutex.Unlock()
	for _, value := range chat.users {
		if value == nickname {
			*reply = "Already user\n"
			return nil	
		}
	}
	chat.users = append(chat.users, nickname)
	chat.messageQueue[nickname] = nil;
	for k, _ := range chat.messageQueue {
		if(k != nickname) {
			chat.messageQueue[k] = append(chat.messageQueue[k], nickname+ " has joined.")
		}
	}
	*reply = "User create with success\n"
	return nil
}
func (chat *ChatServer) ConnectedUsersList(nothing *Nothing, reply *string) error {
	for _, value := range chat.users {
		*reply += value + "\n"
	}
	return nil
}

func (chat *ChatServer) CheckMessages(nickname string, reply *[]string) error {
	chat.mutex.Lock()
	defer chat.mutex.Unlock()
	*reply = chat.messageQueue[nickname]
	chat.messageQueue[nickname] = nil
	return nil
}
func (chat *ChatServer) Quit(nickname string,  reply *string) error {
	chat.mutex.Lock()
	defer chat.mutex.Unlock()
	delete(chat.messageQueue, nickname)
	for i := range chat.users {
		if chat.users[i] == nickname {
			chat.users = append(chat.users[:i], chat.users[i+1:]...)
		}
	}
	for k, v := range chat.messageQueue {
		chat.messageQueue[k] = append(v, nickname+" has logged out.")
	}
	*reply = "User " + nickname + " has logged out."
	return nil
}
func (chat *ChatServer) SendMessage(message Message,  reply *string) error {
	chat.mutex.Lock()
	defer chat.mutex.Unlock()
	for k, _ := range chat.messageQueue {
		chat.messageQueue[k] = append(chat.messageQueue[k], message.Nickname + ": " + message.Text)
	}
	return nil
}

func main() {
	cs := new(ChatServer)
	cs.messageQueue = make(map[string][]string)
	rpc.Register(cs)
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(c)
	}
}
