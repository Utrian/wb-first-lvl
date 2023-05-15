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

func (repo *OrderRepo) TruncateTables() {
	_, err := repo.db.Exec("TRUNCATE items, orders")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("The orders have been cleared.")
}

func (repo *OrderRepo) GetExistingOrder(order_uid string) (models.Order, error) {
	fmt.Println("Получаем items")
	qItems := `
				SELECT
					i.chrt_id, i.track_number, i.price,
					i.rid, i."name", i.sale, i.size,
					i.total_price, i.nm_id, i.brand, i."status"
				FROM items AS i
				WHERE order_id = $1
				`

	rowsI, err := repo.db.Query(qItems, order_uid)
	if err != nil {
		return models.Order{}, err
	}
	defer rowsI.Close()

	itms := make([]models.Item, 0)
	for rowsI.Next() {
		var itm models.Item
		err := rowsI.Scan(
			&itm.ChrtId, &itm.TrackNumber,
			&itm.Price, &itm.Rid, &itm.Name, &itm.Sale,
			&itm.Size, &itm.TotalPrice, &itm.NmId,
			&itm.Brand, &itm.Status,
		)
		if err != nil {
			return models.Order{}, err
		}
		itms = append(itms, itm)
	}

	fmt.Println("Получаем orders")
	qOrder := `
				SELECT *
				FROM orders
				WHERE order_uid = $1
				`
	rowsO, err := repo.db.Query(qOrder, order_uid)
	if err != nil {
		return models.Order{}, err
	}
	defer rowsO.Close()

	ords := make([]models.Order, 0)
	for rowsO.Next() {
		var ord models.Order
		err = rowsO.Scan(
			&ord.OrderUID, &ord.TrackNumber, &ord.Entry,
			&ord.Delivery, &ord.Payment, &ord.Locale,
			&ord.InternalSignature, &ord.CustomerId,
			&ord.DeliveryService, &ord.Shardkey,
			&ord.SmId, &ord.DateCreated, &ord.OofShard,
		)
		if err != nil {
			return models.Order{}, err
		}
		ords = append(ords, ord)
	}
	ord := ords[0]
	ord.Items = itms

	return ord, nil
}

func (repo *OrderRepo) CreateOrder(msg *stan.Msg) {
	fmt.Println("Зашли в CreateOrder")
	order := parse.ParseJsonToOrder(msg)

	fmt.Println("Определяем переменные с запросами")
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

	fmt.Println("Маршализируем delivery и payment")
	jsonDelivery, _ := json.Marshal(order.Delivery)
	jsonPayment, _ := json.Marshal(order.Payment)

	fmt.Println("Отправляем запрос в orders")
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

	fmt.Println("Размещаем items")
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
