package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

var clientsNames   = make(map[net.Conn]string)

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

	go startServer()

	time.Sleep(1 * time.Second)

	serverAddr := "localhost:8080"

	tests := []struct {
		data     string
		expected string
	}{
		{"Hello, Server!", "Acknowledged"},
		{"Test Message", "Acknowledged"},
		{"Another Test", "Acknowledged"},
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
