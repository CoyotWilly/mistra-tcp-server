package main

import (
	"log"
	"math"
	"strconv"
)

func (client *Client) Send(message string) error {
	return client.SendBytes([]byte(message))
}

func (client *Client) SendBytes(message []byte) error {
	_, err := client.connection.Write(message)

	if err != nil {
		err := client.connection.Close()
		if err != nil {
			log.Println("Failed to close connection")

			return err
		}
		client.Server.onClientDisconnect(client, err)
	}

	return err
}

func Dispatch(msg string) float64 {
	value, err := strconv.ParseFloat(msg, TokenValueConversionBase)
	if err != nil {
		log.Println("[ERROR] Failed to parse value:", msg)
		return math.SmallestNonzeroFloat64
	}

	return value
}
