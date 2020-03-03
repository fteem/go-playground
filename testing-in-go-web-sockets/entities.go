package main

import (
	"fmt"
	"time"
)

type Bid struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

type Auction struct {
	ItemID  int   `json:"item_id"`
	EndTime int64 `json:"end_time"`
	Bids    []*Bid
}

func NewAuction(d time.Duration, itemID int, b []*Bid) Auction {
	return Auction{
		ItemID:  itemID,
		EndTime: time.Now().Add(d).Unix(),
		Bids:    b,
	}
}

func (a *Auction) Bid(amount float64, userID int) (*Bid, error) {
	if len(a.Bids) > 0 {
		largestBid := a.Bids[len(a.Bids)-1]
		if largestBid.Amount >= amount {
			return nil, fmt.Errorf("amount must be larger than %.2f", largestBid.Amount)
		}
	}

	if a.EndTime < time.Now().Unix() {
		return nil, fmt.Errorf("auction already closed")
	}

	bid := Bid{
		Amount: amount,
		UserID: userID,
	}

	// mutex lock
	a.Bids = append(a.Bids, &bid)
	// mutex unlock

	return &bid, nil
}
