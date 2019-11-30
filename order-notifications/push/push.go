package push

import (
	"time"

	"github.com/fteem/order-notifications/user"
)

type Notifier struct{}

func (n Notifier) Send(receiver user.User, message string) error {
	// Simulating API call...
	time.Sleep(3 * time.Second)

	return nil
}
