package handlers

import (
	"fmt"
	"io"
	"net"
	"strings"
)

var UserList = CreateUserList()
var MessagePipe = make(chan Request, 100)
var BrodcastPipe = make(chan Request, 100)
var (
	Green = "\033[32m" // Green text color
	Red   = "\033[31m" // Red text color
	White = "\033[37m" // White text color (default)
	Reset = "\033[0m"  // Reset to default terminal color
)

func HandleConnection(conn net.Conn, requests chan<- Request) {

	var message string
	defer conn.Close()
	WelcomeClient(conn)
	name := UserList.GetName(conn)
	message = fmt.Sprintf("[%v] %vUser %v joined chat%v", GetTimestamp(), Green, name, Reset)
	LogWriter(message, LogFile)
	BrodcastPipe <- Request{client: Client{conn: conn, name: name}, data: message + "\n"}

	dataStorage := make([]byte, 1024)

	for {
		n, err := conn.Read(dataStorage)
		if err != nil {
			if err == io.EOF {
				CloseConnection(conn)
				return
			}
			return
		}
		message = strings.TrimSpace(string(dataStorage[:n]))
		requests <- Request{
			client: Client{conn, UserList.GetName(conn)},
			data:   message,
		}
	}

}

func CloseConnection(conn net.Conn) {
	name := UserList.GetName(conn)
	if name != "" {
		UserList.RemoveClient(name)
	}
	message := fmt.Sprintf("[%v] %vUser %v left the chat%v", GetTimestamp(), Red, name, Reset)
	LogWriter(message, LogFile)
	BrodcastPipe <- Request{client: Client{conn: conn, name: name}, data: message + "\n"}
	conn.Close()
}

func (users Users) AddClient(client *Client) {
	Lock.Lock()
	defer Lock.Unlock()
	users[client.name] = client
}

func (users Users) RemoveClient(name string) {
	Lock.Lock()
	defer Lock.Unlock()
	delete(users, name)
}

func (users Users) GetAllClients() ([]*Client, int) {
	Lock.Lock()
	defer Lock.Unlock()
	namesCounter := 0
	var allClients []*Client
	for _, client := range users {
		if client.name != "" {
			namesCounter++
		}
		allClients = append(allClients, client)
	}
	return allClients, namesCounter
}

func (users *Users) NameOccupied(name string) bool {
	allClients, _ := UserList.GetAllClients()
	for _, user := range allClients {
		if user.name == name {
			return true
		}
	}
	return false
}

func (users *Users) GetName(conn net.Conn) string {

	allClients, _ := UserList.GetAllClients()
	for _, user := range allClients {
		if user.conn == conn {
			return user.name
		}
	}
	return ""
}

func CreateUserList() Users {
	return make(map[string]*Client)
}
