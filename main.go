package main

import (
	"log"
	"strings"
)

func main() {
	ApplicationConfiguration.Init()

	connection := StartClientAsync()

	switch strings.ToLower(ApplicationConfiguration.Mode) {
	case "initiator":
		SendInit(<-connection)
		StartServer(connection)
		break
	case "server":
		StartServer(connection)
		break
	default:
		log.Panicln("Invalid mode. Possible options (case insensitive): SERVER | INITIATOR")
	}
}
