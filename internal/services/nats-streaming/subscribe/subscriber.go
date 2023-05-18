package subscribe

import (
	"wb-first-lvl/internal/database/queries"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

const (
	ClusterID = "test-cluster"
	Channel   = "order-notification"
	clientID  = "order-subscriber"
)

type Subscriber struct {
	clusterID string
	clientID  string
	channel   string
	repo      *queries.OrderRepo
	stanConn  stan.Conn
	sub       stan.Subscription
}

func New(repo *queries.OrderRepo) *Subscriber {
	return &Subscriber{
		clusterID: ClusterID,
		clientID:  clientID,
		channel:   Channel,
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
	sc, err := stan.Connect(sb.clusterID, sb.clientID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	sb.stanConn = sc

	return nil
}

func (sb *Subscriber) InitSub() error {
	sub, err := sb.stanConn.Subscribe(sb.channel, sb.repo.CreateOrder)
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
