package main

import "net"

type Client struct {
	conn net.Conn
	name string
}

// Request represents a client request with a timestamp.
type Request struct {
	client Client
	data   []byte
}
