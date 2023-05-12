package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"wb-first-lvl/internal/models"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Настройка пула соединений, с присвоением к глобальной переменной.
// Это позволяет хэндлерам работать с ними.
func init() {
	var err error

	// tools.Load_env()
	// driverName := os.Getenv("DRIVER_NAME")
	// dataSourceName := os.Getenv("DATA_SOURCE_NAME")

	// db, err = sql.Open(driverName, dataSourceName)
	db, err = sql.Open("postgres", "postgres://paul:60880108@172.18.63.119/l0_paul_db")
	if err != nil {
		log.Fatal(err)
	}

	// db.Ping() позволяет проверить что соединение работает.
	// это нужно т.к. sql.Open() не проверяет этого.
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/orders", ordersIndex)
	http.HandleFunc("/orders/show", ordersShow)
	http.ListenAndServe(":3000", nil)
}

func ordersIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		http.Error(
			w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed,
		)
		return
	}
	defer rows.Close()

	ords := make([]*models.Order, 0)
	for rows.Next() {
		ord := new(models.Order)
		err := rows.Scan(
			&ord.OrderUID, &ord.TrackNumber, &ord.Entry,
			&ord.Delivery, &ord.Payment, &ord.Items,
			&ord.Locale, &ord.InternalSignature, &ord.CustomerId,
			&ord.DeliveryService, &ord.Shardkey,
			&ord.SmId, &ord.DateCreated, &ord.OofShard,
		)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		ords = append(ords, ord)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, ord := range ords {
		fmt.Fprint(w, ord)
	}
}

func ordersShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	orderUid := r.FormValue("orderUid")
	if orderUid == "" {
		http.Error(w, http.StatusText(400), 400) // bad request
		return
	}

	row := db.QueryRow("SELECT * FROM ORDERS WHERE order_uid = $1", orderUid)

	ord := new(models.Order)
	err := row.Scan(
		&ord.OrderUID, &ord.TrackNumber, &ord.Entry,
		&ord.Delivery, &ord.Payment, &ord.Items,
		&ord.Locale, &ord.InternalSignature, &ord.CustomerId,
		&ord.DeliveryService, &ord.Shardkey,
		&ord.SmId, &ord.DateCreated, &ord.OofShard,
	)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprint(w, ord)
}
