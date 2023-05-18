package main

import (
	"fmt"
	"os"
	"wb-first-lvl/internal/services/nats-streaming/subscribe"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
)

const (
	clusterID = subscribe.ClusterID
	channel   = subscribe.Channel
	clientID  = "order-publisher"
)

type publisher struct {
	clusterID string
	clientID  string
	channel   string
	stanConn  stan.Conn
}

func newPub() *publisher {
	return &publisher{
		clusterID: clusterID,
		clientID:  clientID,
		channel:   channel,
	}
}

func (p *publisher) initConnection() error {
	sc, err := stan.Connect(p.clusterID, p.clientID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	p.stanConn = sc

	return nil
}

func (p *publisher) cmdPublishJson() {
	for {
		fmt.Printf("Enter the relative or full path to the json file: ")

		var path string
		fmt.Scanln(&path)

		file, err := os.ReadFile(path)
		if err != nil {
			logrus.Info("Incorrect data. Check the path or file. And try again.")
			continue
		}

		err = p.stanConn.Publish(channel, file)
		if err != nil {
			logrus.Info("Failed to send the file. Try again.")
			continue
		}
		logrus.Info("Json has been sent.")
	}
}

func (p *publisher) Close() {
	logrus.Info("Publisher has closed.")
	p.stanConn.Close()
}

func main() {
	defer closer.Close()

	pub := newPub()

	pub.initConnection()
	closer.Bind(pub.Close)

	pub.cmdPublishJson()
}
