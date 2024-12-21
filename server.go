package main

import (
	"crypto/tls"
	"log"
)

type Server struct {
	address            string
	config             *tls.Config
	onClientConnect    func(client *Client)
	onClientDisconnect func(client *Client, err error)
	onMessage          func(client *Client, message string)
}

func New() *Server {
	log.Println("Starting new server instance. Accepting connection at:", ApplicationConfiguration.Server.Binding)
	server := &Server{
		address: ApplicationConfiguration.Server.Binding,
	}

	server.OnClientConnect(func(client *Client) {})
	server.OnMessage(func(client *Client, message string) {})
	server.OnClientDisconnect(func(client *Client, err error) {})

	return server
}

func StartServer() *Server {
	log.Println("Starting server")
	server := &Server{
		address: ApplicationConfiguration.Server.Binding,
	}

	server.OnClientConnect(func(client *Client) {})
	server.OnMessage(func(client *Client, message string) {})
	server.OnClientDisconnect(func(client *Client, err error) {})

	return server
}
