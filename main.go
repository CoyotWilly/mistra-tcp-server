package main

import (
	"flag"
	"log"
	"net"
	"strconv"
)

func main() {
	mode := flag.String("mode", "server", "Application mode")
	address := flag.String("address", "localhost", "Server address")
	port := flag.Int("port", 8080, "Port to run the server on")
	flag.Parse()

	app := *address + ":" + strconv.Itoa(*port)

	if *mode == "server" {
		server := New(app)

		server.OnClientConnect(func(client *Client) {
			log.Println("Client connected.")
		})

		server.OnMessage(func(client *Client, message string) {
			log.Println("Received:", message)
		})

		server.OnClientDisconnect(func(client *Client, err error) {
			log.Println("Disconnected:", client, err)
		})

		server.Listen()
	} else if *mode == "client" {
		log.Println("Running client application. Targeting server at address:", app)

		conn, err := net.Dial("tcp", app)
		if err != nil {
			log.Fatal("Failed to connect to test server")
		}

		_, err = conn.Write([]byte("Test message\n"))
		if err != nil {
			log.Fatal("Failed to send test message.")
		}

		err = conn.Close()
		if err != nil {
			return
		}
	} else {
		log.Panicln("Invalid mode. Possible options: SERVER | CLIENT")
	}
}
