package main

import (
	"log"
	"net"
	"time"
)

type Client struct {
	connection net.Conn
	Server     *Server
}

func StartClientAsync() chan net.Conn {
	done := make(chan net.Conn)
	go func() {
		done <- StartClient()
	}()
	return done
}

func StartClient() net.Conn {
	app := ApplicationConfiguration.Client.Binding
	log.Println("Running client application. Targeting server at address:", app)

	var conn net.Conn
	var err error

	attempt := 1
	for {
		conn, err = net.Dial("tcp", app)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to server (attempt %d/999999), retrying in %v seconds...",
			attempt, ApplicationConfiguration.SleepTime.Seconds())
		attempt++
		time.Sleep(ApplicationConfiguration.SleepTime)
	}

	return conn
}

func SendInit(conn net.Conn) {
	SendMessageFloat(conn, PingToken)
	SendMessageFloat(conn, PongToken)
}
