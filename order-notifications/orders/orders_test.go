package notifications

import "testing"

func TestInformOrderShipped(t *testing.T) {
	user := User{
		Name:  "Peggy",
		Phone: "+12 345 678 999",
	}
	orderID := "12345"

	got := InformOrderShipped(user, orderID)
	want := true

	if want != got {
		t.Errorf("Want '%t', got '%t'", want, got)
	}
}
