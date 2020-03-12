package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func httpToWS(t *testing.T, u string) string {
	t.Helper()

	wsURL, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}

	switch wsURL.Scheme {
	case "http":
		wsURL.Scheme = "ws"
	case "https":
		wsURL.Scheme = "wss"
	}

	return wsURL.String()
}

func newWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	s := httptest.NewServer(h)
	wsURL := httpToWS(t, s.URL)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

func sendMessage(t *testing.T, ws *websocket.Conn, msg inbound) {
	t.Helper()

	m, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ws.WriteMessage(websocket.BinaryMessage, m); err != nil {
		t.Fatalf("%v", err)
	}
}

func receiveWSMessage(t *testing.T, ws *websocket.Conn) outbound {
	t.Helper()

	_, m, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	var reply outbound
	err = json.Unmarshal(m, &reply)
	if err != nil {
		t.Fatal(err)
	}

	return reply
}

func TestBidsHandler(t *testing.T) {
	tcs := []struct {
		name     string
		bids     []*Bid
		duration time.Duration
		message  inbound
		reply    outbound
	}{
		{
			name:     "with good bid",
			bids:     []*Bid{},
			duration: time.Hour * 1,
			message:  inbound{UserID: 1, Amount: 10},
			reply:    outbound{Body: "Bid placed: 10.00"},
		},
		{
			name: "with bad bid",
			bids: []*Bid{
				&Bid{
					UserID: 1,
					Amount: 20,
				},
			},
			duration: time.Hour * 1,
			message:  inbound{UserID: 1, Amount: 10},
			reply:    outbound{Body: "amount must be larger than 20.00"},
		},
		{

			name: "good bid on expired auction",
			bids: []*Bid{
				&Bid{
					UserID: 1,
					Amount: 20,
				},
			},
			duration: time.Hour * -1,
			message:  inbound{UserID: 1, Amount: 30},
			reply:    outbound{Body: "auction already closed"},
		},
	}

	for _, tt := range tcs {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAuction(tt.duration, 1, tt.bids)
			h := bidsHandler{&a}

			s, ws := newWSServer(t, h)
			defer s.Close()
			defer ws.Close()

			sendMessage(t, ws, tt.message)

			reply := receiveWSMessage(t, ws)

			if reply != tt.reply {
				t.Fatalf("Expected '%+v', got '%+v'", tt.reply, reply)
			}
		})
	}
}
