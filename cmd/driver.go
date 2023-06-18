package main

import (
	"fmt"
	"github.com/n1207n/golang-limit-orderbook/entities"
	"github.com/n1207n/golang-limit-orderbook/utils"
	"math/rand"
	"strconv"
	"sync"
)

const (
	N        = 100_000
	MaxPrice = 300
	MinPrice = 150
)

func main() {
	orderbooks := map[string]*entities.OrderBook{
		"aapl": entities.NewOrderBook("aapl"),
		//"hood": entities.NewOrderBook("hood"),
		//"spy":  entities.NewOrderBook("spy"),
		//"amzn": entities.NewOrderBook("amzn"),
		//"meta": entities.NewOrderBook("meta"),
	}

	wg := &sync.WaitGroup{}

	fmt.Println("Populating limit orders")

	// Populate bid orders
	for ticker, ob := range orderbooks {
		wg.Add(1)
		go addTestLimitOrder(wg, ob, ticker, true)
	}

	// Populate ask orders
	for ticker, ob := range orderbooks {
		wg.Add(1)
		go addTestLimitOrder(wg, ob, ticker, false)
	}

	wg.Wait()

	for ticker, ob := range orderbooks {
		fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	}

	fmt.Println("Order matching begins...")
	for ticker, ob := range orderbooks {
		ob.Match()
		fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	}
}

func addTestLimitOrder(wg *sync.WaitGroup, ob *entities.OrderBook, ticker string, isBid bool) {
	for i := 0; i < N; i++ {
		priceString := strconv.FormatFloat(rand.Float64()*(MaxPrice-MinPrice), 'f', 2, 64)

		ob.AddLimitOrder(
			ticker,
			priceString,
			rand.Intn(2000),
			isBid)

		var latestOrder *utils.LimitOrder

		if isBid {
			latestOrder = ob.Bids[0]
		} else {
			latestOrder = ob.Asks[0]
		}

		fmt.Printf("Created limit order - ticker: %s price: %s quantity: %d isBid: %t timestamp: %d \n", ticker, latestOrder.Price.String(), latestOrder.Quantity, latestOrder.IsBid, latestOrder.Timestamp.UnixMilli())
	}

	wg.Done()
}
