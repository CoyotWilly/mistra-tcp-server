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

func ToTokenString(t int) string {
	if t == Init {
		return "Init"
	} else if t == PingToken {
		return "Ping"
	} else if t == PongToken {
		return "Pong"
	}

	return "Unknown"
}

func (t TokenType) String() string {
	switch t {
	case None:
		return "None"
	case Ping:
		return "Ping"
	case Pong:
		return "Pong"
	case Both:
		return "Both"
	}
	return "Unknown"
}
