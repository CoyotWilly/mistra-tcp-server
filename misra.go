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
	fix := false
	switch misra.State {
	case None:
		if value != math.SmallestNonzeroFloat64 {
			misra.Consume(value)
		}
		break
	case Ping:
		log.Println("[INFO] Entering critical section...")
		time.Sleep(ApplicationConfiguration.SleepTime)
		log.Println("[INFO] Leaving critical section...")
		if value != math.SmallestNonzeroFloat64 {
			fix = misra.Consume(value)
		} else {
			misra.Produce(PingToken)
		}
		break
	case Pong:
		misra.Produce(PongToken)
		break
	case Both:
		log.Println("[ERROR] Both tokens held, processing incarnation...")
		misra.Incarnate(misra.Ping)
		misra.Produce(PingToken)
		misra.Produce(PongToken)
		break
	}

	if fix {
		misra.Handle(value)
	}
}

func (misra *Misra) Consume(value float64) bool {
	val := math.Abs(value)
	if math.Abs(misra.Last) > val {
		log.Println("[INFO] Junk received, skipping...")
		return false
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
			break
		case Pong:
			misra.State = Both
			break
		default:
			log.Println("[ERROR] Invalid state for PING token")
			break
		}
	} else if value < 0 {
		misra.Pong = value
		misra.Ping = val

		switch misra.State {
		case None:
			misra.State = Pong
			break
		case Ping:
			misra.State = Both
			break
		default:
			log.Println("[ERROR] Invalid state for PONG token")
			break
		}
	}

	log.Println("[INFO] Consumed token value:", value)
	return misra.State == Both
}

func (misra *Misra) Produce(tokenType int) {
	if tokenType == PingToken {
		SendMessageFloat(misra.Connection, misra.Ping)
		misra.Last = misra.Ping

		if misra.State == Ping {
			misra.State = None
		} else if misra.State == Both {
			misra.State = Pong
		}

	} else if tokenType == PongToken {
		SendMessageFloat(misra.Connection, misra.Pong)
		misra.Last = misra.Pong

		if misra.State == Pong {
			misra.State = None
		} else if misra.State == Both {
			misra.State = Ping
		}
	}

	log.Printf("[INFO] Produced token type: %v\n", tokenType)
}
