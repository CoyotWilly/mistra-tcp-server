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
		log.Println("[INFO] Connection connected")
		if broker == nil {
			broker = <-connection
		}

		if ApplicationConfiguration.Mode == "initiator" {
			log.Println("[INFO] Application lunched in INITIATOR mode")
			misra.State = Both
			misra.Connection = broker
			misra.Produce(PingToken)
			misra.Produce(PongToken)
		}
	})

	server.OnMessage(func(client *Client, message string) {
		log.Println("[INFO] Received:", message)
		if misra.Connection == nil {
			misra.Connection = broker
		}
		misra.Handle(Dispatch(message))
	})

	server.OnClientDisconnect(func(client *Client, err error) {
		log.Println("[INFO] Disconnected:", client, err)
	})

	server.Listen()
}
