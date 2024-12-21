package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	connection net.Conn
	Server     *Server
}

func Connect(connection chan net.Conn, wg *sync.WaitGroup) net.Conn {
	go func() {
		maxAttempts := 5
		interval := 5 * time.Second

		var conn net.Conn
		var err error

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			conn, err = net.Dial("tcp", ApplicationConfiguration.Client.Binding)
			if err == nil {
				break
			}
			log.Printf("Failed to connect to server (attempt %d/%d), retrying in %v...", attempt, maxAttempts, interval)
			time.Sleep(interval)
		}

		if err != nil {
			log.Fatal("Failed to connect to server after multiple attempts")
		} else {
			log.Println("Successfully connected to server. Target: ", ApplicationConfiguration.Client.Binding)
		}

		defer wg.Done()
		connection <- conn
	}()
	return <-connection
}

func DebugClient() {
	log.Println("Running client application. Targeting server at address:", ApplicationConfiguration.Client.Binding)

	conn, err := net.Dial("tcp", ApplicationConfiguration.Client.Binding)
	if err != nil {
		log.Fatal("Failed to connect to server")
	}

	_, err = conn.Write([]byte(ApplicationConfiguration.Message + "\n"))
	if err != nil {
		log.Fatal("Failed to send test message.")
	}

	err = conn.Close()
	if err != nil {
		return
	}
}
