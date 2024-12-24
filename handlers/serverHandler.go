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

func HandleConnection(conn net.Conn, requests chan<- Request) {

	var message string
	defer conn.Close()
	WelcomeClient(conn)
	name := UserList.GetName(conn)
	message = fmt.Sprintf("[%v] User %v joined chat", GetTimestamp(), name)
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
			// fmt.Printf("Error reading data from %s. %v", conn.RemoteAddr(), err) pochemu voznikaet oshibka?
			fmt.Printf("Error reading data from %s.", conn.RemoteAddr())
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
	message := fmt.Sprintf("[%v] User %v left chat", GetTimestamp(), name)
	LogWriter(message, LogFile)
	BrodcastPipe <- Request{client: Client{conn: conn, name: name}, data: message + "\n"}
}

func (m *Users) AddClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[client.name] = client
}

func (m *Users) RemoveClient(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, name)
}

func (m *Users) GetAllClients() ([]*Client, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	namesCounter := 0
	var allClients []*Client
	for _, client := range m.users {
		if client.name != "" {
			namesCounter++
		}
		allClients = append(allClients, client)
	}
	return allClients, namesCounter
}

func (m *Users) NameOccupied(name string) bool {

	allClients, _ := UserList.GetAllClients()
	fmt.Println(len(allClients))
	for _, user := range allClients {
		fmt.Println("Name: " + user.name)
		if user.name == name {
			return true
		}
	}
	return false
}

func (m *Users) GetName(conn net.Conn) string {

	allClients, _ := UserList.GetAllClients()
	for _, user := range allClients {
		if user.conn == conn {
			return user.name
		}
	}
	return ""
}

func CreateUserList() Users {
	return Users{
		users: make(map[string]*Client),
	}
}
