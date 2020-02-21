package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
)

func createUser(t *testing.T, db *gorm.DB) {
	user := User{Username: "jane", Password: "doe123"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Delete(&user)
	})
}

func TestServer(t *testing.T) {
	db, err := gorm.Open("sqlite3", "./cleanups_test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	createUser(t, db)

	r := Router(db)

	tcs := []struct {
		name           string
		username       string
		password       string
		responseBody   string
		responseStatus int
	}{
		{
			name:           "with invalid username-password combo",
			username:       "jane",
			password:       "doe",
			responseBody:   "Forbidden",
			responseStatus: http.StatusForbidden,
		},
		{
			name:           "with valid username-password combo",
			username:       "jane",
			password:       "doe123",
			responseBody:   "6 x 9 = 42",
			responseStatus: http.StatusOK,
		},
	}

	ts := httptest.NewServer(r)
	defer ts.Close()
	client := ts.Client()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/foo", ts.URL), nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req.SetBasicAuth(tc.username, tc.password)

			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			response, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}

			if res.StatusCode != tc.responseStatus {
				t.Errorf("Want '%d', got '%d'", tc.responseStatus, res.StatusCode)
			}

			if strings.TrimSpace(string(response)) != tc.responseBody {
				t.Errorf("Want '%s', got '%s'", tc.responseBody, string(response))
			}
		})
	}
}
