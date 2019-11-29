package sms

import (
	"time"

	"github.com/fteem/order-notifications/user"
)

func Send(receiver user.User, message string) error {
	// Simulating API call...
	time.Sleep(3 * time.Second)

	return nil
}
