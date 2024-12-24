package handlers

import (
	"fmt"
	"io"
	"net"
	"strings"
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
