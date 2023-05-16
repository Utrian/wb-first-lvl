package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"wb-first-lvl/internal/cache"
	"wb-first-lvl/internal/models"
	"wb-first-lvl/internal/services/parse"

	_ "github.com/lib/pq"

	"github.com/nats-io/stan.go"
)

type OrderRepo struct {
	db    *sql.DB
	cache *cache.Cache
}

func NewOrderRepo(db *sql.DB, cache *cache.Cache) *OrderRepo {
	return &OrderRepo{
		db:    db,
		cache: cache,
	}
}

func (repo *OrderRepo) InitCache() error {
	ords, _ := repo.GetAllOrders()
	if err := repo.cache.RestoreCache(&ords); err != nil {
		fmt.Println("The cache was not initialized.")
		return err
	}
	return nil
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
	if ord, b := repo.cache.Get(order_uid); b {
		fmt.Println("This order from cache")
		return ord, nil
	}

	rowsI, err := repo.db.Query(
		`
		SELECT
			i.chrt_id, i.track_number, i.price,
			i.rid, i."name", i.sale, i.size,
			i.total_price, i.nm_id, i.brand, i."status"
		FROM items AS i
		WHERE order_id = $1
		`, order_uid,
	)
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

	// новый запрос в orders

	var ord models.Order
	err = repo.db.QueryRow(
		"SELECT * FROM orders WHERE order_uid = $1", order_uid,
	).Scan(
		&ord.OrderUID, &ord.TrackNumber, &ord.Entry,
		&ord.Delivery, &ord.Payment, &ord.Locale,
		&ord.InternalSignature, &ord.CustomerId,
		&ord.DeliveryService, &ord.Shardkey,
		&ord.SmId, &ord.DateCreated, &ord.OofShard,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return models.Order{}, err
		}
	}

	//

	ord.Items = itms

	return ord, nil
}

func (repo *OrderRepo) GetAllOrders() ([]models.Order, error) {
	rowsO, err := repo.db.Query("SELECT * FROM orders")
	if err != nil {
		return []models.Order{}, err
	}
	defer rowsO.Close()

	countOrds, err := repo.GetOrdersCount()
	if err != nil {
		fmt.Println("There are no orders.")
		return []models.Order{}, err
	}

	ords := make([]models.Order, 0, countOrds)
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
			return []models.Order{}, err
		}

		rowsI, err := repo.db.Query(
			`
			SELECT
					i.chrt_id, i.track_number, i.price,
					i.rid, i."name", i.sale, i.size,
					i.total_price, i.nm_id, i.brand, i."status"
				FROM items AS i
				WHERE order_id = $1
			`, ord.OrderUID,
		)
		if err != nil {
			return []models.Order{}, err
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
				return []models.Order{}, err
			}
			itms = append(itms, itm)
		}
		ord.Items = itms
		ords = append(ords, ord)
	}
	return ords, nil
}

func (repo *OrderRepo) GetOrdersCount() (int, error) {
	var count int
	if err := repo.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			return 0, err
		}
	}
	return count, nil
}

func (repo *OrderRepo) CreateOrder(msg *stan.Msg) {
	ord := parse.ParseJsonToOrder(msg)

	repo.cache.Set(ord.OrderUID, ord, 0)

	jsonDelivery, _ := json.Marshal(ord.Delivery)
	jsonPayment, _ := json.Marshal(ord.Payment)

	_, err := repo.db.Exec(
		`
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
		`,
		ord.OrderUID, ord.TrackNumber, ord.Entry,
		jsonDelivery, jsonPayment, ord.Locale,
		ord.InternalSignature, ord.CustomerId,
		ord.DeliveryService, ord.Shardkey, ord.SmId,
		ord.DateCreated, ord.OofShard,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range ord.Items {
		_, err := repo.db.Exec(
			` 
			INSERT INTO items (
				order_id, chrt_id, track_number,
				price, rid, "name", sale, size,
				total_price, nm_id, brand, "status"
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7,
				$8, $9, $10, $11, $12
			)
			`,
			ord.OrderUID, item.ChrtId, item.TrackNumber,
			item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmId, item.Brand, item.Status,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("Заказ размещен")
}
