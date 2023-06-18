package entities

import (
	"container/heap"
	"fmt"
	"github.com/google/uuid"
	utils "github.com/n1207n/golang-limit-orderbook/utils"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type OrderBook struct {
	ticker string
	Bids   utils.OrderPriorityQueue
	Asks   utils.OrderPriorityQueue
}

// NewOrderBook returns a new instance of OrderBook for a ticker
func NewOrderBook(ticker string) *OrderBook {
	ob := &OrderBook{
		ticker: ticker,
		Bids:   make(utils.OrderPriorityQueue, 0),
		Asks:   make(utils.OrderPriorityQueue, 0),
	}

	heap.Init(&ob.Bids)
	heap.Init(&ob.Asks)
	return ob
}

func intMin(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func (ob *OrderBook) AddLimitOrder(ticker string, priceString string, quantity int, IsBid bool) bool {
	if ticker != ob.ticker {
		log.Fatalf("Wrong ticker to place a new order. Unable to create a new order: %s, %s, %d, IsBid: %t", ticker, priceString, quantity, IsBid)
		return false
	}

	priceDecimal, err := decimal.NewFromString(priceString)

	if err != nil {
		log.Fatalf("Invalid price found. Unable to create a new order: %s, %s, %d, IsBid: %t", ticker, priceString, quantity, IsBid)
		return false
	}

	newOrder := &utils.LimitOrder{
		IsBid:     IsBid,
		ID:        uuid.New(),
		Price:     priceDecimal,
		Quantity:  quantity,
		Ticker:    ticker,
		Timestamp: time.Now(),
	}

	if newOrder.IsBid {
		heap.Push(&ob.Bids, newOrder)
	} else {
		heap.Push(&ob.Asks, newOrder)
	}

	return true
}

func (ob *OrderBook) Match() {
	for ob.Bids.Len() > 0 && ob.Asks.Len() > 0 {
		// Can't fulfill the matching push these orders back to the orderbook
		buy := ob.Bids.Peek().(*utils.LimitOrder)
		sell := ob.Asks.Peek().(*utils.LimitOrder)

		fmt.Printf("bid/ask to match: %d shares at %s VS %d shares at %s\n", buy.Quantity, buy.Price.String(), sell.Quantity, sell.Price.String())

		if buy.Price.LessThan(sell.Price) {
			break
		}

		buy = heap.Pop(&ob.Bids).(*utils.LimitOrder)
		sell = heap.Pop(&ob.Asks).(*utils.LimitOrder)

		quantityFilled := intMin(buy.Quantity, sell.Quantity)
		fmt.Printf("Ticker %s - Matched %d shares at %s\n", ob.ticker, quantityFilled, sell.Price.String())

		buy.Quantity -= quantityFilled
		sell.Quantity -= quantityFilled

		// Order lots are partially fulfilled
		if buy.Quantity > 0 {
			heap.Push(&ob.Bids, buy)
		}

		if sell.Quantity > 0 {
			heap.Push(&ob.Asks, sell)
		}
	}
}
