package main

import (
	"flag"
	"log"
	"strconv"
)

func main() {
	parseFlags()

	if ApplicationConfiguration.Mode == "server" {
		StartMisraServer()
	} else if ApplicationConfiguration.Mode == "client" {
		DebugClient()
	} else {
		log.Panicln("Invalid mode. Possible options: SERVER | CLIENT")
	}
}

func parseFlags() {
	mode := flag.String("mode", "server", "Application mode")
	init := flag.Bool("init", false, "Initialize token passing")
	message := flag.Int("message", 0, "Message to send")

	serverAddress := flag.String("serverAddress", "localhost", "Server address")
	serverPort := flag.Int("serverPort", 5999, "Port to run the server on")

	clientAddress := flag.String("clientAddress", "localhost", "Client data target server IP address")
	clientPort := flag.Int("clientPort", 5998, "Client data target server port")
	flag.Parse()

	app := *serverAddress + ":" + strconv.Itoa(*serverPort)
	clientApp := *clientAddress + ":" + strconv.Itoa(*clientPort)

	ApplicationConfiguration = Configuration{
		Mode:    *mode,
		IsInit:  *init,
		Message: strconv.Itoa(*message),
		Server: Connection{
			Address: *serverAddress,
			Port:    *serverPort,
			Binding: app,
		},
		Client: Connection{
			Address: *clientAddress,
			Port:    *clientPort,
			Binding: clientApp,
		},
	}
}
