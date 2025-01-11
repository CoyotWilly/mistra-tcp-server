package main

import (
	"log"
	"math"
	"net"
	"strings"
	"time"
)

type Misra struct {
	Last       float64
	Ping       float64
	Pong       float64
	State      TokenType
	Connection net.Conn
}

func NewMisra() *Misra {
	if strings.ToLower(ApplicationConfiguration.Mode) == "initiator" {
		return &Misra{Init, PingToken, PongToken, Both, nil}
	}

	return &Misra{Init, PingToken, PongToken, None, nil}
}

func (misra *Misra) Regenerate() {
	log.Printf("[INFO] Regenerating tokens..")
	misra.State = Both
	misra.Ping = math.Abs(misra.Ping) + 1
	misra.Pong = -misra.Ping
	log.Printf("[INFO] Recreated tokens ping: %v | pong: %v\n", misra.Ping, misra.Pong)
}

func (misra *Misra) Incarnate(message float64) {
	log.Println("[ERROR] Incarnate token")
	misra.Ping = math.Abs(message) + 1
	misra.Pong = -misra.Ping
	log.Printf("[INFO] Incarnated token ping: %v | pong: %v\n", misra.Ping, misra.Pong)
}

func (misra *Misra) Handle(value float64) {
	switch misra.State {
	case None:
		if value != math.SmallestNonzeroFloat64 {
			misra.Process(value)
		}
		break
	case Ping:
		log.Println("[INFO] Entering critical section...")
		time.Sleep(ApplicationConfiguration.SleepTime)
		log.Println("[INFO] Leaving critical section...")
		if value != math.SmallestNonzeroFloat64 {
			misra.Process(value)
		} else {
			SendMessageFloat(misra.Connection, PingToken)
		}

	case Pong:
		SendMessageFloat(misra.Connection, PongToken)
	case Both:
		log.Println("[ERROR] Both tokens held, processing incarnation...")
		misra.Incarnate(misra.Ping)
		SendMessageFloat(misra.Connection, PingToken)
		SendMessageFloat(misra.Connection, PongToken)
	}
}

func (misra *Misra) Process(value float64) {
	val := math.Abs(value)
	if math.Abs(misra.Last) > val {
		log.Println("[INFO] Junk received, skipping...")
		return
	}

	if misra.Last == value && value > 0 {
		log.Println("[INFO] Oops something went wrong, PONG has been lost")
		misra.Regenerate()
	} else if misra.Last == value && value < 0 {
		log.Println("[INFO] Oops something went wrong, PING has been lost")
		misra.Regenerate()
	}

	if value > 0 {
		misra.Ping = value
		misra.Pong = -misra.Ping

		switch misra.State {
		case None:
			misra.State = Ping
		case Ping:
			misra.State = Both
		default:
			log.Println("[ERROR] Invalid state for PING token")
		}
	} else if value < 0 {
		misra.Pong = value
		misra.Ping = val

		switch misra.State {
		case None:
			misra.State = Pong
		case Pong:
			misra.State = Both
		default:
			log.Println("[ERROR] Invalid state for PONG token")
		}
	}
}
