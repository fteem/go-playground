package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type AuthMiddleware struct {
	db *gorm.DB
}

func (am *AuthMiddleware) Validate(username, password string) bool {
	var u User

	if err := am.db.Where("username = ? AND password = ?", username, password).First(&u).Error; err != nil {
		return false
	}

	return true
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok && am.Validate(username, password) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func FooHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "6 x 9 = 42\n")
}

func Router(db *gorm.DB) *mux.Router {
	aum := &AuthMiddleware{db}
	r := mux.NewRouter()
	r.Use(aum.Middleware)
	r.HandleFunc("/foo", FooHandler)

	return r
}

func main() {
	db, err := gorm.Open("sqlite3", "./cleanups.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{})

	r := Router(db)

	srvAddress := fmt.Sprintf(":%s", os.Getenv("APP_SERVER_HTTP_PORT"))
	log.Fatal(http.ListenAndServe(srvAddress, r))
}
