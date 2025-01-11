package main

type TokenType int

const (
	None TokenType = iota
	Ping
	Pong
	Both
)

const (
	TokenValueConversionBase = 64
	Init                     = 0
	PingToken                = 1
	PongToken                = -1
)
