package main

import (
	"fmt"
	"os"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-publisher"
	channel   = "order-notification"
)

func InitConnection() stan.Conn {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		logrus.Error(err)
	}
	return sc
}

func cmdPublishJson(sc stan.Conn) {
	for {
		fmt.Print("Enter the full path to the json file: ")

		var path string
		fmt.Scanln(&path)

		if path == "exit" {
			logrus.Info("The channel was shut down.")
			return
		}

		file, err := os.ReadFile(path)
		if err != nil {
			logrus.Info("Incorrect file format or incorrect path. Try again.")
			continue
		}

		err = sc.Publish(channel, file)
		if err != nil {
			logrus.Info("Failed to send the file. Try again.")
			continue
		}
		logrus.Info("Json has been sent.")
	}
}

func main() {
	ic := InitConnection()
	cmdPublishJson(ic)
}
