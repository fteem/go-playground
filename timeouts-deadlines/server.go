package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func slowAPICall(ctx context.Context) string {
	d := rand.Intn(5)
	select {
	case <-ctx.Done():
		log.Printf("slowAPICall was supposed to take %d seconds, but was cancelled.", d)
		return ""
	case <-time.After(time.Duration(d) * time.Second):
		log.Printf("Slow API call done after %d seconds.\n", d)
		return "foobar"
	}
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	result := slowAPICall(r.Context())
	io.WriteString(w, result+"\n")
}

func main() {
	rand.Seed(time.Now().UnixNano())
	srv := http.Server{
		Addr:         ":8888",
		WriteTimeout: 5 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(slowHandler), 1*time.Second, "Timeout!"),
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
