package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "http service address")
)

func main() {
	flag.Parse()

	auction := NewAuction(time.Hour*1, 1, []*Bid{})
	mux := http.NewServeMux()
	mux.Handle("/bids", bidsHandler{&auction})

	log.Fatal(http.ListenAndServe(*addr, mux))
}
