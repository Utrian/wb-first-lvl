package parse

import (
	"encoding/json"
	"wb-first-lvl/internal/models"

	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

func ParseJsonToOrder(orderMsg *stan.Msg) (models.Order, error) {
	var order models.Order
	if err := json.Unmarshal(orderMsg.Data, &order); err != nil {
		logrus.Error(err)
		return order, err
	}
	return order, nil
}
