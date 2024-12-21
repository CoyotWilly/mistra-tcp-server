package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Misra struct {
	Ping      int
	Pong      int
	LastToken int
}

func StartMisraServer() {
	data := Misra{1, -1, 0}
	server := StartServer()
	connection := make(chan net.Conn)

	var wg sync.WaitGroup
	wg.Add(1)

	go Connect(connection, &wg)

	server.OnClientConnect(func(client *Client) {
		log.Println("Client connected.")

		if ApplicationConfiguration.IsInit {
			wg.Wait()

			log.Println("Sending token passing initialization message")
			msg, _ := strconv.Atoi(ApplicationConfiguration.Message)
			send(msg, <-connection)
		}
	})

	server.OnMessage(func(client *Client, message string) {
		log.Println("Received:", strings.TrimRight(message, "\n"))
		data.Receive(message)
		go data.Send(message, <-connection)
	})

	server.OnClientDisconnect(func(client *Client, err error) {
		log.Println("Disconnect")
	})

	log.Println("Accepting connection at:", ApplicationConfiguration.Server.Binding)
	server.Listen()
}

func (misra *Misra) Receive(message string) {
	msg, _ := strconv.Atoi(message)
	if msg == misra.Ping && misra.LastToken == misra.Ping {
		misra.regenerate(msg, "ping")
	} else if msg == misra.Pong && misra.LastToken == misra.Pong {
		misra.regenerate(msg, "pong")
	}

	if misra.Ping == misra.Pong {
		misra.incarnate(msg)
	}
}

func (misra *Misra) Send(message string, connection net.Conn) {
	msg, _ := strconv.Atoi(message)
	if misra.Ping == msg {
		misra.LastToken = misra.Ping
		msg = misra.Ping
	} else if misra.Pong == msg {
		misra.LastToken = misra.Pong
		msg = misra.Pong
	}

	send(msg, connection)
}

func (misra *Misra) regenerate(message int, messageType string) {
	fmt.Printf("[WARN] Regenerating | type: %s | message: %d\n", messageType, message)
	if messageType == "ping" {
		misra.Ping = int(math.Abs(float64(message)))
		misra.Pong = -misra.Ping
	} else if messageType == "pong" {
		misra.Pong = int(math.Abs(float64(message)))
		misra.Ping = -misra.Pong
	}
}

func (misra *Misra) incarnate(message int) {
	misra.Ping = int(math.Abs(float64(message))) + 1
	misra.Pong = -misra.Ping
}

func send(message int, connection net.Conn) {
	_, err := connection.Write([]byte(strconv.Itoa(message) + "\n"))
	if err != nil {
		log.Fatal("Failed to send message")
	}
}
