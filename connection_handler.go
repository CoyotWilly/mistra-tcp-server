package main

import "net"

func (client *Client) Connection() net.Conn {
	return client.connection
}

func (client *Client) Close() error {
	return client.connection.Close()
}
