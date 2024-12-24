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

func ClientName(client net.Conn) {
	_, namesCounter := UserList.GetAllClients()
	if namesCounter > RoomSize-1 {
		client.Write([]byte("\nSorry, chat room is full, please try again later.\n"))
		client.Close()
		return
	}

	client.Write([]byte("[ENTER YOUR NAME]: "))
	scanner := bufio.NewScanner(client)

	for {
		scanner.Scan()
		name := strings.TrimSpace(scanner.Text())
		switch nameIsValid(name) {
		case 0:
			client.Write([]byte("Name should contain visible characters. Please enter a valid name.\n[ENTER YOUR NAME]: "))
		case -1:
			client.Write([]byte("Name contains invalid characters. Please choose another name.\n[ENTER YOUR NAME]: "))
		default:
			if !UserList.NameOccupied(name) {
				UserList.AddClient(&Client{
					conn: client,
					name: name,
				})
				client.Write([]byte(name + ", welcome to our chat!\nType " + Red + "/exit" + Reset + " when you want to leave the chat.\n"))
				// send history to user

				return
			} else {
				client.Write([]byte("Name is already taken. Please choose another name.\n[ENTER YOUR NAME]: "))
			}
		}

	}

}

func nameIsValid(name string) int {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return 0
	}
	for _, ch := range name {
		if ch < 32 || ch > 126 {
			return -1
		}
	}
	return 1

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
