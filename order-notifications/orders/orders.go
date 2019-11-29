package notifications

import (
	"fmt"

	"github.com/fteem/order-notifications/sms"
	"github.com/fteem/order-notifications/user"
)

func InformOrderShipped(receiver user.User, orderID string) bool {
	message := fmt.Sprintf("Your order #%s is shipped!", orderID)
	err := sms.Send(receiver, message)

	if err != nil {
		return false
	}

	return true
}
