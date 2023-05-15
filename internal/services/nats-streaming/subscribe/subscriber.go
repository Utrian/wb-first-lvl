package subscribe

import (
	"fmt"
	"wb-first-lvl/internal/database/queries"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	ClusterID string
	ClientID  string
	Channel   string
	repo      queries.OrderRepo
}

func New(repo queries.OrderRepo) *Subscriber {
	return &Subscriber{
		ClusterID: "test-cluster",
		ClientID:  "order-subscriber",
		Channel:   "order-nitofication",
		repo:      repo,
	}
}

func (sb *Subscriber) SubAndPub() {
	sc, err := stan.Connect(sb.ClusterID, sb.ClientID)
	if err != nil {
		fmt.Println(err)
	}
	defer sc.Close()

	fmt.Println("Подписываемся на стриминг")
	sub, err := sc.Subscribe(sb.Channel, sb.repo.CreateOrder, stan.StartWithLastReceived())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Вышли из стриминга")

	defer sub.Unsubscribe()
}

// func (sb *Subscriber) CreateOrder(m *stan.Msg) {
// 	order := models.Order{}

// 	err := json.Unmarshal(m.Data, &order)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	if err := sb.repo.CreateOrder(&order); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

// func (sb *Subscriber) ConnAndCreate() {
// 	sc, err := stan.Connect(sb.ClusterID, sb.ClientID)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}
// 	defer sc.Close()

// 	sub, err := sc.Subscribe(sb.Channel, sb.CreateOrder, stan.StartWithLastReceived())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	time.Sleep(5 * time.Second)
// 	defer sub.Unsubscribe()
// }
