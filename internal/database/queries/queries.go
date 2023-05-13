package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"

	"wb-first-lvl/internal/models"
	"wb-first-lvl/tools"
)

func GetAllOrders() {
	tools.Load_env()
	driverName := os.Getenv("DRIVER_NAME")
	dataSourceName := os.Getenv("DATA_SOURCE_NAME")

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	ords := make([]*models.Order, 0)
	for rows.Next() {
		ord := new(models.Order)
		err := rows.Scan(
			&ord.OrderUID, &ord.TrackNumber,
			&ord.Entry, &ord.Delivery, &ord.Payment,
			&ord.Items, &ord.Locale, &ord.InternalSignature,
			&ord.CustomerId, &ord.DeliveryService, &ord.Shardkey,
			&ord.SmId, &ord.DateCreated, &ord.OofShard,
		)
		if err != nil {
			log.Fatal(err)
		}
		ords = append(ords, ord)
	}
	if err = rows.Err(); err != nil { // Нужен ли этот блок?
		log.Fatal(err)
	}

	for _, ord := range ords {
		fmt.Println(ord) // Возможно потребуется форматировать вывод!
	}
}

func CreateOrder(orderMsg *stan.Msg) {
	// Парсим json файл в структуру
	var order models.Order
	if err := json.Unmarshal(orderMsg.Data, &order); err != nil {
		fmt.Println(err)
	}
	fmt.Println(order)

	// Подключаемся к БД и выполняем запрос на вставку
	// На данный момент нужно разобраться с стркутурой []item
	tools.Load_env()
	driverName := os.Getenv("DRIVER_NAME")
	dataSourceName := os.Getenv("DATA_SOURCE_NAME")

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jsonDelivery, _ := json.Marshal(order.Delivery)
	jsonPayment, _ := json.Marshal(order.Payment)
	jsonItems, _ := json.Marshal(order.Items)

	_, err = db.Exec(
		`
		INSERT INTO orders (
			order_uid, track_number, "entry",
			delivery, payment, items, locale,
			internal_signature, customer_id,
			delivery_service, shardkey, sm_id,
			date_created, off_shard
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14
		)
		`,
		order.OrderUID, order.TrackNumber, order.Entry,
		jsonDelivery, jsonPayment, jsonItems, order.Locale,
		order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.Shardkey, order.SmId,
		order.DateCreated, order.OofShard,
	)
	if err != nil {
		panic(err)
	}
}
