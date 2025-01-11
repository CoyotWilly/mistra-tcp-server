package main

import (
	"flag"
	"strconv"
	"time"
)

var ApplicationConfiguration Configuration

type Configuration struct {
	Mode      string
	SleepTime time.Duration
	Server    Connection
	Client    Connection
}

type Connection struct {
	Address string
	Port    int
	Binding string
}

func (config Configuration) Init() {
	mode := flag.String("mode", "server", "Application mode")
	sleep := flag.Float64("sleep", 5, "Number of seconds to sleep between sending messages")

	serverAddress := flag.String("serverAddress", "localhost", "Server listen address")
	serverPort := flag.Int("serverPort", 5999, "Server listen port")

	clientAddress := flag.String("clientAddress", "localhost", "Destination server address - the address of the destination server where all the message will be sent")
	clientPort := flag.Int("clientPort", 5998, "Destination server port")

	flag.Parse()

	server := *serverAddress + ":" + strconv.Itoa(*serverPort)
	app := *clientAddress + ":" + strconv.Itoa(*clientPort)

	ApplicationConfiguration = Configuration{
		Mode:      *mode,
		SleepTime: time.Duration(*sleep) * time.Second,
		Server: Connection{
			Address: *serverAddress,
			Port:    *serverPort,
			Binding: server,
		},
		Client: Connection{
			Address: *clientAddress,
			Port:    *clientPort,
			Binding: app,
		},
	}
}
