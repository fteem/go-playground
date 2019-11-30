package orders

import (
	"fmt"

	"github.com/fteem/order-notifications/user"
)

type Sender interface {
	Send(user.User, string) error
}

func InformOrderShipped(receiver user.User, orderID string, sender Sender) bool {
	message := fmt.Sprintf("Your order #%s is shipped!", orderID)
	err := sender.Send(receiver, message)

	if err != nil {
		return false
	}

	return true
}
