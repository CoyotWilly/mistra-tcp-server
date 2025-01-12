package main

import (
	"log"
	"strings"
	"time"
)

func main() {
	ApplicationConfiguration.Init()

	connection := StartClientAsync()
	mode := strings.ToLower(ApplicationConfiguration.Mode)

	if mode == "initiator" || mode == "server" {
		StartServer(connection)
	} else if mode == "client" {
		conn := <-connection
		for {
			SendMessage(conn, "Misra test only client sample message!")
			time.Sleep(1 * time.Second)
		}
	} else {
		log.Panicln("Invalid mode. Possible options (case insensitive): SERVER | CLIENT | INITIATOR")
	}
}
