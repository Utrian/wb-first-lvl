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
	repo      *queries.OrderRepo
	stanConn  stan.Conn
	sub       stan.Subscription
}

func New(repo *queries.OrderRepo) *Subscriber {
	return &Subscriber{
		ClusterID: "test-cluster",
		ClientID:  "order-subscriber",
		Channel:   "order-notification",
		repo:      repo,
	}
}

func (sb *Subscriber) SubAndPub() error {
	if err := sb.InitConn(); err != nil {
		return err
	}

	if err := sb.InitSub(); err != nil {
		return err
	}

	select {}
}

func (sb *Subscriber) InitConn() error {
	sc, err := stan.Connect(sb.ClusterID, sb.ClientID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	sb.stanConn = sc

	return nil
}

func (sb *Subscriber) InitSub() error {
	sub, err := sb.stanConn.Subscribe(sb.Channel, sb.repo.CreateOrder)
	if err != nil {
		logrus.Error(err)
		return err
	}
	sb.sub = sub

	return nil
}

func (sb *Subscriber) Close() {
	logrus.Info("Nats-streaming has unsubscribed and closed.")
	sb.sub.Unsubscribe()
	sb.stanConn.Close()
}
