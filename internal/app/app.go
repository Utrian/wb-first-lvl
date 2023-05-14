package app

import rec "wb-first-lvl/internal/services/nats-streaming/receive"

func Run() {
	rec.Receiver()
}
