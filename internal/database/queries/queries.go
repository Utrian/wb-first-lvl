package queries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"wb-first-lvl/internal/cache"
	"wb-first-lvl/internal/models"
	"wb-first-lvl/internal/services/parse"
	"wb-first-lvl/tools"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/nats-io/stan.go"
)

type OrderRepo struct {
	db    *sql.DB
	cache *cache.Cache
}

func NewOrderRepo() *OrderRepo {
	db, err := InitConn()
	if err != nil {
		logrus.Error("The database is not available.")
		return &OrderRepo{}
	}

	cache := cache.NewCache(5*time.Minute, 10*time.Minute)

	return &OrderRepo{
		db:    db,
		cache: cache,
	}
}

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

func (repo *OrderRepo) InitCache() error {
	ords, err := repo.GetAllOrders()
	if err != nil {
		logrus.Error("The cache was not initialized.")
		return err
	}

	repo.cache.RestoreCache(&ords)

	return nil
}

func (repo *OrderRepo) TruncateTables() {
	_, err := repo.db.Exec("TRUNCATE items, orders")
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("The orders have been cleared.")
}

func (repo *OrderRepo) GetExistingOrder(order_uid string) (models.Order, error) {
	if ord, b := repo.cache.Get(order_uid); b {
		infoMsg := fmt.Sprintf("Order %s from cache.", order_uid)
		logrus.Info(infoMsg)
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
		logrus.Error(err)
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
			logrus.Error(err)
			return models.Order{}, err
		}
		itms = append(itms, itm)
	}

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
			errMsg := fmt.Sprintf("The order number %s does not exist.", order_uid)
			logrus.Info(errMsg)
			return models.Order{}, err
		}
	}

	ord.Items = itms

	infoMsg := fmt.Sprintf("Order %s from database.", order_uid)
	logrus.Info(infoMsg)

	return ord, nil
}

func (repo *OrderRepo) GetAllOrders() ([]models.Order, error) {
	rowsO, err := repo.db.Query("SELECT * FROM orders")
	if err != nil {
		logrus.Error(err)
		return []models.Order{}, err
	}
	defer rowsO.Close()

	countOrds, err := repo.GetOrdersCount()
	if err != nil {
		logrus.Info("There are no orders.")
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
			logrus.Error(err)
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
			logrus.Error(err)
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
				logrus.Error(err)
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
			logrus.Error(err)
			return 0, err
		}
	}
	return count, nil
}

func (repo *OrderRepo) CreateOrder(msg *stan.Msg) {
	ord, err := parse.ParseJsonToOrder(msg)
	if err != nil {
		return
	}

	if v := ord.Validator(); !v {
		return
	}

	repo.cache.Set(ord.OrderUID, ord, 0)

	jsonDelivery, _ := json.Marshal(ord.Delivery)
	jsonPayment, _ := json.Marshal(ord.Payment)

	_, err = repo.db.Exec(
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
		logrus.Error(err)
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
			logrus.Error(err)
			return
		}
	}
	logrus.Info("Заказ размещен")
}

func (repo *OrderRepo) Close() {
	logrus.Info("DB has closed.")
	repo.db.Close()
}
