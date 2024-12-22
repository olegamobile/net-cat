package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func handleConnection(conn net.Conn, requests chan<- Request) {

	var message string
	defer conn.Close()
	WelcomeClient(conn)
	fmt.Printf("Client connected (%s)\n", conn.RemoteAddr())

	dataStorage := make([]byte, 1024)

	for {
		n, err := conn.Read(dataStorage)
		if err != nil {
			if err == io.EOF {
				CloseConnection(conn)
				return
			}
			fmt.Printf("Error reading data from %s. %v", conn.RemoteAddr(), err)
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
	fmt.Printf("%v left the chat %s\n", name, conn.RemoteAddr())
}

func processRequests(requests <-chan Request) {
	for request := range requests {
		if len(request.data) != 0 {
			now := time.Now()
			timestamp := now.Format("2006-01-02 15:04:05")
			fmt.Printf("[%v][%v]: %s\n", timestamp, request.client.name, request.data)
		}
	}
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

func (m *Users) GetAllClients() []*Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	var allClients []*Client
	for _, client := range m.users {
		allClients = append(allClients, client)
	}
	return allClients
}

func (m *Users) NameOccupied(name string) bool {

	allClients := UserList.GetAllClients()
	for _, user := range allClients {
		fmt.Println("Name: " + user.name)
		if user.name == name {
			return true
		}
	}
	return false
}

func (m *Users) GetName(conn net.Conn) string {

	allClients := UserList.GetAllClients()
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
