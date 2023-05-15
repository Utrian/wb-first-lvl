package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"wb-first-lvl/tools"

	"wb-first-lvl/internal/database/queries"
	"wb-first-lvl/internal/services/nats-streaming/subscribe"
	// rec "wb-first-lvl/internal/services/nats-streaming/receive"
)

func InitConn() *sql.DB {
	tools.Load_env()
	driverName := os.Getenv("DRIVER_NAME")
	dataSourceName := os.Getenv("DATA_SOURCE_NAME")

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return db
}

func Run() {
	// Подключаемся к БД и создаем соединение repo
	db := InitConn()
	defer db.Close()
	repo := queries.NewOrderRepo(db)
	repo.TruncateTables()

	// Подключаемся к стриммингу и делаем запись в БД
	sub := subscribe.New(*repo)
	sub1 := *sub.SubAndPub()
	defer sub1.Unsubscribe()
	time.Sleep(200 * time.Millisecond)
}
