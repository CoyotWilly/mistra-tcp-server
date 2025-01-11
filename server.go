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

func StartServerAsync(connection chan net.Conn) chan bool {
	done := make(chan bool)
	go func() {
		StartServer(connection)
		done <- true
	}()
	return done
}

func StartServer(connection chan net.Conn) {
	misra := NewMisra()
	server := New(ApplicationConfiguration.Server.Binding)

	server.OnClientConnect(func(client *Client) {
		log.Println("Connection connected.")
	})

	server.OnMessage(func(client *Client, message string) {
		log.Println("Received:", message)
		misra.Connection = <-connection
		misra.Handle(Dispatch(message))
	})

	server.OnClientDisconnect(func(client *Client, err error) {
		log.Println("Disconnected:", client, err)
	})

	server.Listen()
}
