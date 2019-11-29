package notifications

import (
	"errors"
	"testing"

	"github.com/fteem/order-notifications/user"
)

func TestInformOrderShipped(t *testing.T) {
	cases := []struct {
		user         user.User
		orderID      string
		sendingError error
		name         string
		want         bool
	}{
		{
			user:         user.User{"Peggy", "+12 345 678 999"},
			orderID:      "12345",
			sendingError: nil,
			want:         true,
			name:         "Successful send",
		},
		{
			user:         user.User{"Peggy", "+12 345 678 999"},
			orderID:      "12345",
			sendingError: errors.New("Sending failed"),
			want:         false,
			name:         "Unsuccessful send",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockSend := func(user.User, string) error {
				return tc.sendingError
			}

			got := InformOrderShipped(tc.user, tc.orderID, mockSend)

			if tc.want != got {
				t.Errorf("Want '%t', got '%t'", tc.want, got)
			}
		})
	}
}
