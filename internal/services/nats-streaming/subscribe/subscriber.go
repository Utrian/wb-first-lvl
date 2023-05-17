package subscribe

import (
	"wb-first-lvl/internal/database/queries"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type Subscriber struct {
	ClusterID string
	ClientID  string
	Channel   string
	repo      queries.OrderRepo
}

func New(repo queries.OrderRepo) *Subscriber {
	return &Subscriber{
		ClusterID: "test-cluster",
		ClientID:  "order-subscriber",
		Channel:   "order-notification",
		repo:      repo,
	}
}

func (sb *Subscriber) SubAndPub() *stan.Subscription {
	sc, err := stan.Connect(sb.ClusterID, sb.ClientID)
	if err != nil {
		logrus.Error(err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe(sb.Channel, sb.repo.CreateOrder)
	if err != nil {
		logrus.Error(err)
	}

	defer sub.Unsubscribe()

	select {} // позволяет функции дальше слушать канал и сразу обрабатывать по поступлению json;
	// надо найти как правильно завершать процесс в таком случае;
}
