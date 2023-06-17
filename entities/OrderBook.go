package entities

import (
	"container/heap"
	"fmt"
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

func intMin(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func (ob *OrderBook) Match() {
	for ob.bids.Len() > 0 && ob.asks.Len() > 0 {
		buy := heap.Pop(&ob.bids).(*LimitOrder)
		sell := heap.Pop(&ob.asks).(*LimitOrder)

		// Can't fulfill the matching push these orders back to the orderbook
		if buy.Price < sell.Price {
			heap.Push(&ob.bids, buy)
			heap.Push(&ob.asks, sell)
			break
		}

		quantity_filled := intMin(buy.Quantity, sell.Quantity)
		fmt.Printf("Matched %d shares at %d\n", quantity_filled, sell.Price)

		buy.Quantity -= quantity_filled
		sell.Quantity -= quantity_filled

		// Order lots are partially fulfilled
		if buy.Quantity > 0 {
			heap.Push(&ob.bids, buy)
		}

		if sell.Quantity > 0 {
			heap.Push(&ob.asks, sell)
		}
	}
}
