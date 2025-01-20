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
	go misra.tryInitiate(connection)

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

func (misra *Misra) tryInitiate(broker chan net.Conn) {
	if ApplicationConfiguration.Mode != "initiator" {
		return
	}

	for {
		log.Println("[INFO] Waiting for connection...")
		con := <-broker
		if con != nil {
			log.Println("[INFO] Application lunched in INITIATOR mode")
			misra.State = Both
			misra.Connection = con
			misra.Produce(PingToken)
			misra.Produce(PongToken)
			break
		}
	}
}
