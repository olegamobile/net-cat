package handlers

import (
	"bufio"
	"net"
	"strings"
)

func CreateUserList() Users {
	return make(map[string]*Client)
}

func WelcomeClient(client net.Conn) {
	client.Write([]byte("Welcome to TCP-Chat!" + welcomeIcon))
	ClientName(client)

	client.Write([]byte(MsgHistory))
}

func ClientName(conn net.Conn) {
	namesCounter := len(UserList)
	if namesCounter > RoomSize-1 {
		conn.Write([]byte("\nSorry, chat room is full, please try again later.\n"))
		conn.Close()
		return
	}

	conn.Write([]byte("[ENTER YOUR NAME]: "))
	scanner := bufio.NewScanner(conn)

	for {
		scanner.Scan()
		name := clearInput(scanner.Text())
		if nameIsValid(conn, name) {
			if !UserList.NameOccupied(name) {
				UserList.AddClient(&Client{
					conn: conn,
					name: name,
				})
				conn.Write([]byte(name + ", welcome to our chat!\nType " + Red + "/exit" + Reset + " when you want to leave the chat.\n"))
				return
			} else {
				conn.Write([]byte("Name is already taken. Please choose another name.\n[ENTER YOUR NAME]: "))
			}
		}
	}
}

func nameIsValid(conn net.Conn, name string) bool {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		conn.Write([]byte("Name should contain visible characters. Please enter a valid name.\n[ENTER YOUR NAME]: "))
		return false
	}
	for _, ch := range name {
		if ch < 32 || ch > 126 {
			conn.Write([]byte("Name contains invalid characters. Please choose another name.\n[ENTER YOUR NAME]: "))
			return false
		}
	}
	return true

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

func (users Users) GetAllClients() []*Client {
	Lock.Lock()
	defer Lock.Unlock()
	var allClients []*Client
	for _, client := range users {
		allClients = append(allClients, client)
	}
	return allClients
}

func (users *Users) NameOccupied(name string) bool {
	allClients := UserList.GetAllClients()
	for _, user := range allClients {
		if user.name == name {
			return true
		}
	}
	return false
}

func (users *Users) GetName(conn net.Conn) string {

	allClients := UserList.GetAllClients()
	for _, user := range allClients {
		if user.conn == conn {
			return user.name
		}
	}
	return ""
}
