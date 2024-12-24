package main

import (
	"fmt"
	"net"
	"net_cat/handlers"
	"os"
)

func main() {
	port := getPort()
	StartServer(port)
}

func StartServer(port string) {
	listener, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	handlers.LogFileCreate()
	defer handlers.LogFile.Close()

	go handlers.ProcessMessages(handlers.MessagePipe)
	go handlers.BroadcastMessages(handlers.BrodcastPipe)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handlers.HandleConnection(conn, handlers.MessagePipe)

	}

}

func getPort() string {
	args := os.Args[1:]
	argsLen := len(args)
	switch {
	case argsLen == 0:
		return "8989"
	case argsLen == 1:
		return args[0]
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}
	return ""
}
