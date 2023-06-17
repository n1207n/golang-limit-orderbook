package entities

import (
	utils "github.com/n1207n/golang-limit-orderbook/utils"
)

type OrderBook struct {
	ticker string
	bids   utils.PriorityQueue
	asks   utils.PriorityQueue
}

// NewOrderBook returns a new instance of OrderBook for a ticker
func NewOrderBook(ticker string) *OrderBook {
	return &OrderBook{
		ticker: ticker,
		bids:   make(utils.PriorityQueue, 0),
		asks:   make(utils.PriorityQueue, 0),
	}
}
