package main

func (server *Server) OnClientConnect(callback func(client *Client)) {
	server.onClientConnect = callback
}

func (server *Server) OnClientDisconnect(callback func(client *Client, err error)) {
	server.onClientDisconnect = callback
}

func (server *Server) OnMessage(callback func(client *Client, message string)) {
	server.onMessage = callback
}
