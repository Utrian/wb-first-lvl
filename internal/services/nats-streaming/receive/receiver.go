package receive

import (
	"fmt"
	"wb-first-lvl/internal/database/queries"

	"github.com/nats-io/stan.go"
)

const (
	ClusterID = "test-cluster"
	ClientID  = "order-subscriber"
	Channel   = "order-notification"
)

func Receiver() {
	sc, err := stan.Connect(ClusterID, ClientID)
	if err != nil {
		fmt.Println(err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe(Channel, queries.CreateOrder, stan.StartWithLastReceived())
	if err != nil {
		fmt.Println(err)
	}
	defer sub.Unsubscribe()
}
