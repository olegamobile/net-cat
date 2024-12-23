package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

var clientsNames = make(map[net.Conn]string)

func testClientSendData(t *testing.T, serverAddr string, data string) string {

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(data))
	if err != nil {
		t.Fatalf("Failed to write to server: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from server: %v", err)
	}

	return string(buf[:n])
}

func TestServerCommunication(t *testing.T) {
	port := "8989"
	go StartServer(port)

	time.Sleep(1 * time.Second)

	serverAddr := "localhost:" + port

	tests := []struct {
		data     string
		expected string
	}{
		{"Hello, Server!", `Welcome to TCP-Chat!
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
	"     `-'       `--'\n" + "[ENTER YOUR NAME]: "},
	}

	// running test
	for _, test := range tests {
		t.Run(fmt.Sprintf("Sending: %s", test.data), func(t *testing.T) {
			// imitating sending data and receiving a response
			response := testClientSendData(t, serverAddr, test.data)

			// checking that the server response matches the expected
			if response != test.expected {
				t.Errorf("Expected response %s, but got %s", test.expected, response)
			}
		})
	}
}
