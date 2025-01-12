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
	var broker net.Conn
	misra := NewMisra()
	server := NewServer()

	server.OnClientConnect(func(client *Client) {
		log.Println("[INFO] Connection connected.")
		if broker == nil {
			broker = <-connection
		}

		if ApplicationConfiguration.Mode == "initiator" {
			SendInit(broker)
		}
	})

	server.OnMessage(func(client *Client, message string) {
		log.Println("[INFO] Received:", message)
		SendMessage(broker, message)
		misra.Connection = broker
		misra.Handle(Dispatch(message))
	})

	server.OnClientDisconnect(func(client *Client, err error) {
		log.Println("[INFO] Disconnected:", client, err)
	})

	server.Listen()
}
