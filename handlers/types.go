package handlers

import (
	"net"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

// Request represents a client request with a timestamp.
type Request struct {
	client Client
	data   string
}

type Users struct {
	users map[string]*Client
	mu    sync.Mutex
}

const RoomSize = 10
const LogFileDir = "logs"
var Port = "8989"