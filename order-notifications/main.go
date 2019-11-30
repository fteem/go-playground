package main

import (
	"github.com/fteem/order-notifications/orders"
	"github.com/fteem/order-notifications/push"
	"github.com/fteem/order-notifications/sms"
	"github.com/fteem/order-notifications/user"
)

func main() {
	u := user.User{"Peggy", "+123 456 789"}
	orderID := "123"
	dispatcher := sms.Dispatcher{}
	notifier := push.Notifier{}
	orders.InformOrderShipped(u, orderID, dispatcher)
	orders.InformOrderShipped(u, orderID, notifier)
}
