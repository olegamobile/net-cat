package main

import (
	"fmt"
	"net"
)



func main() {
	UserList = CreateUserList()
	startServer()
}

func startServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on localhost:8080...")

	requests := make(chan Request, 100)

	go processRequests(requests)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, requests)

	}

}
