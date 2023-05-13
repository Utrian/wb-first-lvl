package app

import (
	"wb-first-lvl/internal/database/queries"
	rec "wb-first-lvl/internal/services/nats-streaming/receive"
)

func Run() {
	rec.Receiver()
	queries.GetAllOrders()
}
