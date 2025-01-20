package main

import (
	"crypto/tls"
	"log"
	"net"
)

type Server struct {
	address            string
	config             *tls.Config
	onClientConnect    func(client *Client)
	onClientDisconnect func(client *Client, err error)
	onMessage          func(client *Client, message string)
}

func StartServer(connection chan net.Conn) {
	misra := NewMisra()
	server := NewServer()
	go misra.TryInitiate(connection)

	server.OnClientConnect(func(client *Client) {
		log.Println("[INFO] Client connection established")
	})

	server.OnMessage(func(client *Client, message string) {
		log.Println("[INFO] Received:", message)
		if misra.Connection == nil {
			misra.Connection = <-connection
		}
		misra.Handle(Dispatch(message))
	})

	server.OnClientDisconnect(func(client *Client, err error) {
		log.Println("[INFO] Disconnected:", client, err)
	})

	server.Listen()
}
