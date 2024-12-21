package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
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

func handleConnection(conn net.Conn, requests chan<- Request) {

	defer conn.Close()
	WelcomeClient(conn)
	fmt.Printf("Client connected (%s)\n", conn.RemoteAddr())

	dataStorage := make([]byte, 1024)

	for { 
		_, err := conn.Read(dataStorage)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client left the chat %s\n", conn.RemoteAddr())
				return
			}
			fmt.Printf("Error reading data from %s. %v", conn.RemoteAddr(), err)
			return
		}

		// fmt.Println("[" + string(dataStorage) + "]")
		requests <- Request{
			client: Client{conn, ""},
			data:   dataStorage,
		}
	}

}

func processRequests(requests <-chan Request) {
	for request := range requests {
		fmt.Printf("Received: %s", string(request.data))
	}
}
