package main

import "net"

type Client struct {
	connection net.Conn
	Server     *Server
}
