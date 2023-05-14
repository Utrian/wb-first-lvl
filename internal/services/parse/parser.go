package parse

import (
	"encoding/json"
	"fmt"
	"wb-first-lvl/internal/models"

	"github.com/nats-io/stan.go"
)

func ParseJsonToOrder(orderMsg *stan.Msg) models.Order {
	var order models.Order
	if err := json.Unmarshal(orderMsg.Data, &order); err != nil {
		fmt.Println(err)
	}
	return order
}
