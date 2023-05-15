package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"wb-first-lvl/internal/models"
	"wb-first-lvl/internal/services/parse"

	_ "github.com/lib/pq"

	"github.com/nats-io/stan.go"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (repo *OrderRepo) GetExistingOrder(order_uid string) (models.Order, error) {
	rows, err := repo.db.Query("SELECT * FROM orders WHERE order_uid = $1", order_uid)
	if err != nil {
		fmt.Println(err)
		return models.Order{}, err
	}
	defer rows.Close()

	ord := new(models.Order)
	err = rows.Scan(
		&ord.OrderUID, &ord.TrackNumber,
		&ord.Entry, &ord.Delivery, &ord.Payment,
		&ord.Items, &ord.Locale, &ord.InternalSignature,
		&ord.CustomerId, &ord.DeliveryService, &ord.Shardkey,
		&ord.SmId, &ord.DateCreated, &ord.OofShard,
	)
	if err != nil {
		fmt.Println(err)
		return *ord, err
	}
	if err = rows.Err(); err != nil { // Нужен ли этот блок?
		fmt.Println(err)
		return *ord, err
	}

	return *ord, nil
}

func (repo *OrderRepo) CreateOrder(msg *stan.Msg) {
	order := parse.ParseJsonToOrder(msg)
	fmt.Println(order)

	fmt.Println("Начинаем запрос")
	var (
		qOrder = `
			INSERT INTO orders (
				order_uid, track_number, "entry",
				delivery, payment, locale,
				internal_signature, customer_id,
				delivery_service, shardkey, sm_id,
				date_created, off_shard
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7,
				$8, $9, $10, $11, $12, $13
			)
		`
		qItems = `
			INSERT INTO items (
				order_id, chrt_id, track_number,
				price, rid, "name", sale, size,
				total_price, nm_id, brand, "status"
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7,
				$8, $9, $10, $11, $12
			)
		`
	)

	jsonDelivery, _ := json.Marshal(order.Delivery)
	jsonPayment, _ := json.Marshal(order.Payment)

	_, err := repo.db.Exec(
		qOrder,
		order.OrderUID, order.TrackNumber, order.Entry,
		jsonDelivery, jsonPayment, order.Locale,
		order.InternalSignature, order.CustomerId,
		order.DeliveryService, order.Shardkey, order.SmId,
		order.DateCreated, order.OofShard,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Заказ размещен")

	for _, item := range order.Items {
		_, err := repo.db.Exec(
			qItems,
			order.OrderUID, item.ChrtId, item.TrackNumber,
			item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmId, item.Brand, item.Status,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Айтемы размещены")
}

// func (repo *OrderRepo) CreateOrder(order *models.Order) error {
// 	fmt.Println("Начинаем запрос")
// 	var (
// 		qOrder = `
// 			INSERT INTO orders (
// 				order_uid, track_number, "entry",
// 				delivery, payment, locale,
// 				internal_signature, customer_id,
// 				delivery_service, shardkey, sm_id,
// 				date_created, off_shard
// 			) VALUES (
// 				$1, $2, $3, $4, $5, $6, $7,
// 				$8, $9, $10, $11, $12, $13
// 			)
// 		`
// 		qItems = `
// 			INSERT INTO items (
// 				order_id, chrt_id, track_number,
// 				price, rid, "name", sale, size,
// 				total_price, nm_id, brand, "status"
// 			) VALUES (
// 				$1, $2, $3, $4, $5, $6, $7,
// 				$8, $9, $10, $11, $12
// 			)
// 		`
// 	)

// 	jsonDelivery, _ := json.Marshal(order.Delivery)
// 	jsonPayment, _ := json.Marshal(order.Payment)

// 	_, err := repo.db.Exec(
// 		qOrder,
// 		order.OrderUID, order.TrackNumber, order.Entry,
// 		jsonDelivery, jsonPayment, order.Locale,
// 		order.InternalSignature, order.CustomerId,
// 		order.DeliveryService, order.Shardkey, order.SmId,
// 		order.DateCreated, order.OofShard,
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println("Заказ размещен")

// 	for _, item := range order.Items {
// 		_, err := repo.db.Exec(
// 			qItems,
// 			order.OrderUID, item.ChrtId, item.TrackNumber,
// 			item.Price, item.Rid, item.Name, item.Sale, item.Size,
// 			item.TotalPrice, item.NmId, item.Brand, item.Status,
// 		)
// 		if err != nil {
// 			fmt.Println(err)
// 			return err
// 		}
// 	}
// 	fmt.Println("Айтемы размещены")

// 	return nil
// }
