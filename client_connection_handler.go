package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

func (client *Client) listen() {
	client.Server.onClientConnect(client)
	reader := bufio.NewReader(client.connection)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			err := client.connection.Close()
			if err != nil {
				return
			}

			client.Server.onClientDisconnect(client, err)

			return
		}

		client.Server.onMessage(client, message)
	}
}

func (server *Server) Listen() {
	var listener net.Listener
	var err error

	if server.config == nil {
		listener, err = net.Listen("tcp", server.address)
	} else {
		listener, err = tls.Listen("tcp", server.address, server.config)
	}

	if err != nil {
		log.Fatalf("Error starting TCP server on %s with %s\n", server.address, err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Error closing TCP listener %s\n", err)
		}
	}(listener)

	for {
		conn, _ := listener.Accept()
		client := &Client{
			connection: conn,
			Server:     server,
		}

		go client.listen()
	}
}
