package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const welcomeIcon = `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
|    ` + "`.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n"

func WelcomeClient(client net.Conn) {

	client.Write([]byte("Welcome to TCP-Chat!" + welcomeIcon))
	ClientName(client)
}

func ClientName(client net.Conn) {

	client.Write([]byte("[ENTER YOUR NAME]: "))
	scanner := bufio.NewScanner(client)

	for {
		scanner.Scan()
		name := scanner.Text()
		switch nameIsValid(name) {
		case 0:
			client.Write([]byte("Name should contain visible characters. Please enter a valid name.\n[ENTER YOUR NAME]: "))
		case -1:
			client.Write([]byte("Name contains invalid characters. Please choose another name.\n[ENTER YOUR NAME]: "))
		default:
			fmt.Println(UserList.NameOccupied(name))
			if !UserList.NameOccupied(name) {
				UserList.AddClient(&Client{
					conn: client,
					name: name,
				})
				client.Write([]byte(name + ", welcome to our chat!\n"))
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
