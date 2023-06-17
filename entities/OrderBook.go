package entities

import (
	"fmt"
	"github.com/google/uuid"
	utils "github.com/n1207n/golang-limit-orderbook/utils"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type OrderBook struct {
	ticker string
	bids   utils.OrderPriorityQueue
	asks   utils.OrderPriorityQueue
}

// NewOrderBook returns a new instance of OrderBook for a ticker
func NewOrderBook(ticker string) *OrderBook {
	return &OrderBook{
		ticker: ticker,
		bids:   make(utils.OrderPriorityQueue, 0),
		asks:   make(utils.OrderPriorityQueue, 0),
	}
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
		ob.bids.Push(newOrder)
	} else {
		ob.asks.Push(newOrder)
	}

	return true
}

func (ob *OrderBook) Match() {
	for ob.bids.Len() > 0 && ob.asks.Len() > 0 {
		buy := ob.bids[0]
		sell := ob.asks[0]

		// Can't fulfill the matching push these orders back to the orderbook
		if buy.Price.LessThan(sell.Price) {
			break
		}

		buy = ob.bids.Pop().(*utils.LimitOrder)
		sell = ob.asks.Pop().(*utils.LimitOrder)

		quantity_filled := intMin(buy.Quantity, sell.Quantity)
		fmt.Printf("Matched %d shares at %d\n", quantity_filled, sell.Price)

		buy.Quantity -= quantity_filled
		sell.Quantity -= quantity_filled

		// Order lots are partially fulfilled
		if buy.Quantity > 0 {
			ob.bids.Push(buy)
		}

		if sell.Quantity > 0 {
			ob.asks.Push(sell)
		}
	}
}
