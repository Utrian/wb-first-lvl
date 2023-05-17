package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"wb-first-lvl/internal/cache"
	"wb-first-lvl/internal/database/queries"
	"wb-first-lvl/internal/transport/router"
	"wb-first-lvl/tools"
	// rec "wb-first-lvl/internal/services/nats-streaming/receive"
)

func InitConn() (*sql.DB, error) {
	tools.Load_env()
	driverName := os.Getenv("DRIVER_NAME")
	dataSourceName := os.Getenv("DATA_SOURCE_NAME")

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Run() {
	// Подключаемся к БД и создаем соединение repo
	db, err := InitConn()
	if err != nil {
		fmt.Println("The database is not available.")
		return
	}
	defer db.Close()

	cache := cache.NewCache(5*time.Minute, 10*time.Minute)
	repo := queries.NewOrderRepo(db, cache)
	repo.InitCache()

	router.Router(repo)

	// sub := subscribe.New(*repo)
	// sub.SubAndPub()

	// Проверка на то что из кеша возвращаются данные
	// ord, _ := repo.GetExistingOrder("b563feb7b2b84b6test")
	// fmt.Printf("detailed struct: %+v/n", ord)

	// Подключаемся к стриммингу и делаем запись в БД
	// sub := subscribe.New(*repo)
	// go sub.SubAndPub()
	// time.Sleep(200 * time.Millisecond)
}
