package main

import (
	"log"
	"math"
	"math/rand"
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
	misra.Ping = math.Abs(misra.Ping)
	misra.Pong = -misra.Ping
	misra.State = Both
	log.Printf("[INFO] Recreated tokens ping: %v | pong: %v | state: %s\n", misra.Ping, misra.Pong, misra.State.String())
}

func (misra *Misra) Incarnate(message float64) {
	log.Println("[ERROR] Incarnate token")
	misra.Ping = math.Abs(message) + 1
	misra.Pong = -misra.Ping
	log.Printf("[INFO] Incarnated token ping: %v | pong: %v | state: %s\n", misra.Ping, misra.Pong, misra.State.String())
}

func (misra *Misra) Handle(value float64) {
	repeat := true

	for repeat {
		switch misra.State {
		case None:
			if value != math.SmallestNonzeroFloat64 {
				misra.Consume(value)
				repeat = false
			}
			break
		case Ping:
			if rand.Float64() > ApplicationConfiguration.LoseProbability {
				log.Println("[INFO] Losing token...")
				misra.Last = misra.Ping

				if misra.State == Ping {
					misra.State = None
				} else if misra.State == Both {
					misra.State = Pong
				}
				break
			}

			log.Println("[INFO] Entering critical section...")
			time.Sleep(ApplicationConfiguration.SleepTime)
			log.Println("[INFO] Leaving critical section...")
			if value != math.SmallestNonzeroFloat64 {
				misra.Consume(value)
				repeat = misra.State == Both || misra.State == Pong
			} else {
				misra.Produce(PingToken)
			}
			break
		case Pong:
			misra.Produce(PongToken)
			repeat = false
			break
		case Both:
			log.Println("[ERROR] Both tokens held, processing incarnation...")
			misra.Incarnate(misra.Ping)
			misra.Produce(PingToken)
			misra.Produce(PongToken)
			repeat = false
			break
		}
	}
}

func (misra *Misra) Consume(value float64) {
	val := math.Abs(value)
	if math.Abs(misra.Last) > val {
		log.Println("[INFO] Junk received, skipping...")
		return
	}

	if misra.Last == value && misra.Last > Init {
		log.Println("[INFO] Oops something went wrong, PONG has been lost")
		misra.Regenerate()
	} else if misra.Last == value && misra.Last < Init {
		log.Println("[INFO] Oops something went wrong, PING has been lost")
		misra.Regenerate()
	}

	if value > Init {
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
	} else if value < Init {
		misra.Ping = val
		misra.Pong = value

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

	log.Printf("[INFO] Consumed token value: %.0f | type %s\n", value, misra.State.String())
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

	log.Printf("[INFO] Produced token type: %s\n", ToTokenString(tokenType))
}
