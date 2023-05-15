package main

import (
	"fmt"
	"os"

	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "order-publisher"
	channel   = "order-notification"
)

func InitConnection() stan.Conn {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		panic(err)
	}
	return sc
}

func cmdPublishJson(sc stan.Conn) {
	for {
		fmt.Print("Enter the full path to the json file: ")

		var path string
		fmt.Scanln(&path)

		if path == "exit" {
			fmt.Println("The channel was shut down.")
			return
		}

		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Incorrect file format or incorrect path. Try again.")
			continue
		}

		err = sc.Publish(channel, file)
		if err != nil {
			fmt.Println("Failed to send the file. Try again.")
			continue
		}
		fmt.Println("Json has been sent.")
	}
}

func main() {
	ic := InitConnection()
	cmdPublishJson(ic)
}
