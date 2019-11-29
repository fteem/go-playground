package notifications

import (
	"fmt"

	"github.com/fteem/order-notifications/user"
)

func InformOrderShipped(receiver user.User, orderID string, sendSMS func(user.User, string) error) bool {
	message := fmt.Sprintf("Your order #%s is shipped!", orderID)
	err := sendSMS(receiver, message)

	if err != nil {
		return false
	}

	return true
}
