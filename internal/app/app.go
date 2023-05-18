package app

import (
	"wb-first-lvl/internal/database/queries"
	"wb-first-lvl/internal/services/nats-streaming/subscribe"
	"wb-first-lvl/internal/transport/router"

	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
)

func Run() {
	defer closer.Close()
	closer.Bind(accessClose)

	repo := queries.NewOrderRepo()
	repo.InitCache()

	closer.Bind(repo.Close)

	sub := subscribe.New(repo)
	go sub.SubAndPub()

	closer.Bind(sub.Close)

	router.Router(repo)
}

func accessClose() {
	logrus.Info("The program has been closed.")
}
