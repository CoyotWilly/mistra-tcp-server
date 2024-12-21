package main

var ApplicationConfiguration Configuration

type Configuration struct {
	Mode    string
	IsInit  bool
	Message string
	Server  Connection
	Client  Connection
}

type Connection struct {
	Address string
	Port    int
	Binding string
}
