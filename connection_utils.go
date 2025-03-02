package main

import (
	"fmt"
	"log"
	"net"
)

func SendMessageFloat(connection net.Conn, message float64) {
	SendMessage(connection, fmt.Sprintf("%.0f", message))
}

func SendMessage(connection net.Conn, message string) {
	log.Printf("[INFO] Sending message {string} %s | Target: {address: %s }", message, connection.RemoteAddr().String())
	_, err := connection.Write([]byte(message + "\n"))
	if err != nil {
		log.Fatal("[ERROR] Failed to send message", err)
	}
}
