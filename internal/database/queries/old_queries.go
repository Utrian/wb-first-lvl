package queries

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"

// 	_ "github.com/lib/pq"
// 	"github.com/nats-io/stan.go"

// 	"wb-first-lvl/internal/models"
// 	"wb-first-lvl/internal/services/parse"
// 	"wb-first-lvl/tools"
// )

// func InitConn() *sql.DB {
// 	tools.Load_env()
// 	driverName := os.Getenv("DRIVER_NAME")
// 	dataSourceName := os.Getenv("DATA_SOURCE_NAME")

// 	db, err := sql.Open(driverName, dataSourceName)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	return db
// }

// func GetAllOrders() {
// 	db := InitConn()
// 	defer db.Close()

// 	rows, err := db.Query("SELECT * FROM orders")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	ords := make([]*models.Order, 0)
// 	for rows.Next() {
// 		ord := new(models.Order)
// 		err := rows.Scan(
// 			&ord.OrderUID, &ord.TrackNumber,
// 			&ord.Entry, &ord.Delivery, &ord.Payment,
// 			&ord.Items, &ord.Locale, &ord.InternalSignature,
// 			&ord.CustomerId, &ord.DeliveryService, &ord.Shardkey,
// 			&ord.SmId, &ord.DateCreated, &ord.OofShard,
// 		)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		ords = append(ords, ord)
// 	}
// 	if err = rows.Err(); err != nil { // Нужен ли этот блок?
// 		log.Fatal(err)
// 	}

// 	for i, ord := range ords {
// 		fmt.Println(i, ord) // Возможно потребуется форматировать вывод!
// 	}
// }

// func OrderExists(db *sql.DB, uid string) bool {
// 	fmt.Println("Зашли в функцию проверки")
// 	fmt.Println("Начинаем проверку")
// 	order, err := db.Exec("SELECT * FROM orders WHERE order_uid = $1", uid)
// 	fmt.Println("Проверка завершена")
// 	if err != nil {
// 		panic(err)
// 	}
// 	isExist, _ := order.RowsAffected()
// 	fmt.Println(isExist == 0)
// 	return isExist == 1
// }

// func CreateOrder(orderMsg *stan.Msg) {
// 	db := InitConn()
// 	defer db.Close()

// 	fmt.Println("Началась запись в БД")
// 	order := parse.ParseJsonToOrder(orderMsg)
// 	fmt.Println(order)

// 	fmt.Println("Проверяем наличие заказа")
// 	if orderExists := OrderExists(db, order.OrderUID); orderExists {
// 		fmt.Println("This order already exists.")
// 		return
// 	}

// 	fmt.Println("Преобразуем структуры в json")
// 	// Преобразуем структуры в json, чтобы сохранить их в виде jsonb в БД.
// 	jsonDelivery, _ := json.Marshal(order.Delivery)
// 	jsonPayment, _ := json.Marshal(order.Payment)
// 	jsonItems, _ := json.Marshal(order.Items)

// 	fmt.Println("Началась запись в БД")
// 	_, err := db.Exec(
// 		`
// 		INSERT INTO orders (
// 			order_uid, track_number, "entry",
// 			delivery, payment, items, locale,
// 			internal_signature, customer_id,
// 			delivery_service, shardkey, sm_id,
// 			date_created, off_shard
// 		) VALUES (
// 			$1, $2, $3, $4, $5, $6, $7,
// 			$8, $9, $10, $11, $12, $13, $14
// 		)
// 		`,
// 		order.OrderUID, order.TrackNumber, order.Entry,
// 		jsonDelivery, jsonPayment, jsonItems, order.Locale,
// 		order.InternalSignature, order.CustomerId,
// 		order.DeliveryService, order.Shardkey, order.SmId,
// 		order.DateCreated, order.OofShard,
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	fmt.Println("The order has been successfully created.")
// }

// func CreateTables() {
// 	db := InitConn()
// 	defer db.Close()

// 	query, err := ioutil.ReadFile("internal/database/postgres.sql")
// 	if err != nil {
// 		panic(err)
// 	}

// 	_, err = db.Exec(string(query))
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func TruncateOrders() {
// 	db := InitConn()
// 	defer db.Close()

// 	_, err := db.Exec("TRUNCATE orders")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("The orders have been cleared.")
// }

// func DropTables() {
// 	db := InitConn()
// 	defer db.Close()

// 	_, err := db.Exec("DROP TABLE orders, items")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("The orders have been deleted.")
// }
