package main

import (
	"fmt"
	"net/rpc"
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	nickname string
	Connection   *rpc.Client
}
type Message struct {
	Nickname   string
	Text string
}
var wg sync.WaitGroup

const (
	CMD_PREFIX = "/"
	CMD_CREATE = CMD_PREFIX + "create"
	CMD_LIST   = CMD_PREFIX + "list"
	CMD_HELP   = CMD_PREFIX + "help"
	CMD_QUIT   = CMD_PREFIX + "quit"

	MSG_HELP = "\nCommands:\n" +
		CMD_CREATE + " nickname - create a user named nickame\n" +
		CMD_LIST + " - lists all users\n" +
		CMD_HELP + " - lists all commands\n" +
		CMD_QUIT + " - quits the program\n" + 
		"write your message\n"
	MSG_DISCONNECT = "Disconnected from the server.\n"
)
func (c *Client) CreateConnection() {
	connection, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatalln(err)
		return
	}
	c.Connection = connection
}
func (c *Client) CheckMessages() {
	var reply []string
	for {
		err := c.Connection.Call("ChatServer.CheckMessages", c.nickname, &reply)
		if err != nil {
			log.Fatalln("Chat has been shutdown. Goodbye.")
		}
		for i := range reply {
			log.Println(reply[i])
		}
		time.Sleep(time.Second)
	}
}

func (c *Client) Input() {
	for {
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if err != nil {
			wg.Done()
			break
		}
		if strings.HasPrefix(str, CMD_CREATE) {
			c.CreateUser(str)
		} else if strings.HasPrefix(str, CMD_LIST) {
			c.ConnectedUsersList() 
		} else if strings.HasPrefix(str, CMD_QUIT) {
			c.Quit() 
		} else if strings.HasPrefix(str, CMD_HELP) {
			c.Help() 
		} else if len(str) > 1 && len(c.nickname) > 0  {
			c.SendMessage(str)	
		} else if( len(c.nickname) == 0) {
			fmt.Println("Create a user with ==> " + CMD_CREATE)
		}
	}
}


func (c *Client) CreateUser(str string) {
	var message string		
	nickname := strings.TrimSuffix(strings.TrimPrefix(str, CMD_CREATE + " "), "\n")
	c.nickname = nickname
	err := c.Connection.Call("ChatServer.CreateUser", nickname, &message)
	if err != nil {
		wg.Done()
	}
	fmt.Print(message)
}
func (c *Client) ConnectedUsersList() {
	var message string		
	err := c.Connection.Call("ChatServer.ConnectedUsersList", true, &message)
	if err != nil {
		wg.Done()
	}
	fmt.Print(message)
}

func (c *Client) Quit() {
	var message string		
	err := c.Connection.Call("ChatServer.Quit", c.nickname, &message)
	if err != nil {
		wg.Done()
	}
	fmt.Print(message)
}
func (c *Client) Help() {
	fmt.Println(MSG_HELP)
}
func (c *Client) SendMessage(str string) {
	var message string			
	text := Message{
		Nickname: c.nickname,
		Text: str,
	}
	err := c.Connection.Call("ChatServer.SendMessage", text, &message)
	if err != nil {
		wg.Done()
	}
	fmt.Print(message)
}



func main() {
	var client *Client = &Client{}
	client.CreateConnection()
	client.Help()
	go client.CheckMessages()
	client.Input()
}
