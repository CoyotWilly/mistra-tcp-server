package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"strings"
)

func (client *Client) listen() {
	client.Server.onClientConnect(client)
	reader := bufio.NewReader(client.connection)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("[ERROR] Couldn't read message:", message, err)
			err := client.connection.Close()
			if err != nil {
				log.Println("[ERROR] Couldn't close connection:", err)
				return
			}

			client.Server.onClientDisconnect(client, err)

			return
		}

		client.Server.onMessage(client, strings.TrimSpace(message))
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
		log.Fatalf("[ERROR] Error starting TCP server on %s with %s\n", server.address, err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("[ERROR] Error closing TCP listener %s\n", err)
		}
	}(listener)
	log.Printf("[INFO] Server ready to accept connections on %s\n", server.address)

	for {
		conn, _ := listener.Accept()
		client := &Client{
			connection: conn,
			Server:     server,
		}

		log.Println("[INFO] New connection from", conn.RemoteAddr().String())
		go client.listen()
	}
}

func NewServer() *Server {
	address := ApplicationConfiguration.Server.Binding
	log.Println("[INFO] Starting new server instance. Accepting connection at:", address)
	server := &Server{
		address: address,
	}

	server.OnClientConnect(func(client *Client) {})
	server.OnMessage(func(client *Client, message string) {})
	server.OnClientDisconnect(func(client *Client, err error) {})

	return server
}
