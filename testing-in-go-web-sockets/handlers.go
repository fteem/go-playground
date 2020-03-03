package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type inbound struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

type outbound struct {
	Body string `json:"body"`
}

type bidsHandler struct {
	auction *Auction
}

func (bh bidsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		var in inbound
		err = json.Unmarshal(m, &in)
		if err != nil {
			handleError(ws, err)
			continue
		}

		bid, err := bh.auction.Bid(in.Amount, in.UserID)
		if err != nil {
			handleError(ws, err)
			continue
		}

		out, err := json.Marshal(outbound{Body: fmt.Sprintf("Bid placed: %.2f", bid.Amount)})
		if err != nil {
			handleError(ws, err)
			continue
		}

		err = ws.WriteMessage(websocket.BinaryMessage, out)
		if err != nil {
			handleError(ws, err)
			continue
		}
	}
}

func handleError(ws *websocket.Conn, err error) {
	log.Println("Error:", err)

	b, err := json.Marshal(&outbound{Body: err.Error()})
	if err != nil {
		log.Println("Error:", err)
	}

	err = ws.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Println("Error:", err)
	}
}
