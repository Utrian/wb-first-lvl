package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"wb-first-lvl/internal/cache"
	"wb-first-lvl/internal/database/queries"
	"wb-first-lvl/tools"
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

	cache := cache.NewCache(5*time.Minute, 10*time.Minute)
	repo := queries.NewOrderRepo(db, cache)
	repo.InitCache()

	// Проверка на то что из кеша возвращаются данные
	ords, _ := repo.GetExistingOrder("b563feb7b2b84b6test")
	fmt.Println(ords)

	// ord, _ := repo.GetExistingOrder("b563feb7b2b84b6test")
	// fmt.Println("Order: ", ord)

	// Подключаемся к стриммингу и делаем запись в БД
	// sub := subscribe.New(*repo)
	// sub.SubAndPub()
	time.Sleep(200 * time.Millisecond)
}
